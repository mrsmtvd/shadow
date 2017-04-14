package metrics

import (
	"fmt"
	"sync"
	"time"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/snitch"
	_ "github.com/kihamo/snitch/collector"
)

const (
	ComponentName = "metrics"

/*
	TagAppName    = "app_name"
	TagAppVersion = "app_version"
	TagAppBuild   = "app_build"
	TagHostname   = "hostname"
*/
)

type Component struct {
	application shadow.Application

	config *config.Component
	logger logger.Logger

	mutex        sync.RWMutex
	prefix       string
	changeTicker chan time.Duration
}

type hasMetrics interface {
	Metrics() snitch.Collector
}

func (c *Component) GetName() string {
	return ComponentName
}

func (c *Component) GetVersion() string {
	return ComponentVersion
}

func (c *Component) GetDependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: logger.ComponentName,
		},
		{
			Name: dashboard.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(*config.Component)
	c.changeTicker = make(chan time.Duration)

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) error {
	c.logger = logger.NewOrNop(c.GetName(), c.application)

	url := c.config.GetString(ConfigMetricsUrl)
	if url == "" {
		wg.Done()
		return fmt.Errorf("%s is empty", ConfigMetricsUrl)
	}

	c.prefix = c.config.GetString(ConfigMetricsPrefix)

	// search metrics
	components, err := c.application.GetComponents()
	if err != nil {
		return err
	}

	for _, component := range components {
		if metrics, ok := component.(hasMetrics); ok {
			snitch.DefaultRegisterer.Register(metrics.Metrics())
		}
	}

	return nil
}

/*
func (c *Component) initLabels(labels string) {
	l := metric.Labels{
		TagAppName:    c.application.GetName(),
		TagAppVersion: c.application.GetVersion(),
		TagAppBuild:   c.application.GetBuild(),
	}

	if hostname, err := os.Hostname(); err == nil {
		l[TagHostname] = hostname
	}

	if len(labels) > 0 {
		var parts []string

		for _, tag := range strings.Split(labels, ",") {
			parts = strings.Split(tag, "=")

			if len(parts) > 1 {
				l[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.labels = l
}

func (c *Component) GetLastUpdated() *time.Time {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.updatedAt
}
*/
