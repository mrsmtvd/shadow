package internal

import (
	"os"
	"strings"
	"sync"

	"github.com/kihamo/snitch"
	"github.com/kihamo/snitch/collector"
	"github.com/kihamo/snitch/storage"
	"github.com/mrsmtvd/shadow"
	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/i18n"
	"github.com/mrsmtvd/shadow/components/logging"
	"github.com/mrsmtvd/shadow/components/metrics"
	"github.com/mrsmtvd/shadow/components/profiling"
)

const (
	TagAppName    = "app_name"
	TagAppVersion = "app_version"
	TagAppBuild   = "app_build"
	TagHostname   = "hostname"

	StorageIDInflux = "influx"
)

type Component struct {
	application shadow.Application

	config config.Component
	logger logging.Logger

	mutex    sync.RWMutex
	registry snitch.Registerer
}

func (c *Component) Name() string {
	return metrics.ComponentName
}

func (c *Component) Version() string {
	return metrics.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: dashboard.ComponentName,
		},
		{
			Name: i18n.ComponentName,
		},
		{
			Name: logging.ComponentName,
		},
		{
			Name: profiling.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a

	c.config = a.GetComponent(config.ComponentName).(config.Component)

	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) error {
	c.logger = logging.DefaultLazyLogger(c.Name())

	<-a.ReadyComponent(config.ComponentName)
	c.registry = snitch.NewRegistry(c.config.Duration(metrics.ConfigInterval))

	if err := c.initStorage(); err != nil {
		return err
	}

	if c.application.HasComponent(profiling.ComponentName) {
		c.registry.AddStorages(storage.NewExpvarWithID(metrics.ComponentName))
	}

	c.registry.Register(
		NewOnOffCollector(collector.NewDebugCollector(), func() bool {
			return c.config.Bool(metrics.ConfigCollectorsDebugEnabled)
		}),
		NewOnOffCollector(collector.NewRuntimeCollector(), func() bool {
			return c.config.Bool(metrics.ConfigCollectorsRuntimeEnabled)
		}),
		NewOnOffCollector(collector.NewModCollector(), func() bool {
			return c.config.Bool(metrics.ConfigCollectorsModEnabled)
		}),
	)

	c.initLabels(c.config.String(metrics.ConfigLabels))

	ready <- struct{}{}

	// search metrics
	components, err := c.application.GetComponents()
	if err != nil {
		return err
	}

	for _, component := range components {
		if m, ok := component.(metrics.HasMetrics); ok {
			c.Register(m.Metrics())
		}
	}

	return nil
}

func (c *Component) Shutdown() error {
	_, err := c.registry.Gather()
	return err
}

func (c *Component) Registry() snitch.Registerer {
	return c.registry
}

func (c *Component) Register(cs ...snitch.Collector) {
	c.registry.Register(cs...)
}

func (c *Component) initStorage() (err error) {
	url := c.config.String(metrics.ConfigURL)
	if url == "" {
		return
	}

	existsStorage, err := c.registry.GetStorage(StorageIDInflux)
	if err != nil {
		newStorage, err := storage.NewInfluxWithID(
			StorageIDInflux,
			url,
			c.config.String(metrics.ConfigDatabase),
			c.config.String(metrics.ConfigUsername),
			c.config.String(metrics.ConfigPassword),
			c.config.String(metrics.ConfigPrecision))

		if err != nil {
			return err
		}

		c.registry.AddStorages(newStorage)

		return nil
	}

	if castStorage, ok := existsStorage.(*storage.Influx); ok {
		err = castStorage.Reinitialization(
			url,
			c.config.String(metrics.ConfigDatabase),
			c.config.String(metrics.ConfigUsername),
			c.config.String(metrics.ConfigPassword),
			c.config.String(metrics.ConfigPrecision))
	}

	return err
}

func (c *Component) initLabels(labels string) {
	l := snitch.Labels{
		&snitch.Label{Key: TagAppName, Value: c.application.Name()},
		&snitch.Label{Key: TagAppVersion, Value: c.application.Version()},
		&snitch.Label{Key: TagAppBuild, Value: c.application.Build()},
	}

	if hostname, err := os.Hostname(); err == nil {
		l = append(l, &snitch.Label{
			Key:   TagHostname,
			Value: hostname,
		})
	}

	if len(labels) > 0 {
		var parts []string

		for _, tag := range strings.Split(labels, ",") {
			parts = strings.Split(tag, "=")

			if len(parts) > 1 {
				l = append(l, &snitch.Label{
					Key:   strings.TrimSpace(parts[0]),
					Value: strings.TrimSpace(parts[1]),
				})
			}
		}
	}

	c.mutex.Lock()
	c.registry.SetLabels(l)
	c.mutex.Unlock()
}
