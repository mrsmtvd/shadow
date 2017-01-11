package metrics

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	kit "github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/influx"
	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
)

const (
	TagAppName    = "app_name"
	TagAppVersion = "app_version"
	TagAppBuild   = "app_build"
	TagHostname   = "hostname"
)

type Component struct {
	application shadow.Application

	config *config.Component
	logger logger.Logger

	mutex        sync.RWMutex
	client       influxdb.Client
	connector    *influx.Influx
	prefix       string
	changeTicker chan time.Duration
}

type hasMetrics interface {
	MetricsRegister(*Component)
}

type hasCapture interface {
	MetricsCapture()
}

func (c *Component) GetName() string {
	return "metrics"
}

func (c *Component) GetVersion() string {
	return "1.0.0"
}

func (c *Component) Init(a shadow.Application) error {
	resourceConfig, err := a.GetComponent("config")
	if err != nil {
		return err
	}
	c.config = resourceConfig.(*config.Component)

	c.application = a

	c.changeTicker = make(chan time.Duration)

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) error {
	if err := c.initClient(c.config.GetString(ConfigMetricsUrl), c.config.GetString(ConfigMetricsUsername), c.config.GetString(ConfigMetricsPassword)); err != nil {
		wg.Done()
		return err
	}

	c.logger = logger.NewOrNop(c.GetName(), c.application)

	c.connector = influx.New(c.getTags(), influxdb.BatchPointsConfig{
		Database:  c.config.GetString(ConfigMetricsDatabase),
		Precision: c.config.GetString(ConfigMetricsPrecision),
	}, logger.NewGoKitLogger(c.logger))

	c.prefix = c.config.GetString(ConfigMetricsPrefix)

	// search metrics
	for _, component := range c.application.GetComponents() {
		if metrics, ok := component.(hasMetrics); ok {
			metrics.MetricsRegister(c)
		}
	}

	// send to influx
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(c.config.GetDuration(ConfigMetricsInterval))

		for {
			select {
			case <-ticker.C:
				// collect metrics
				wg := new(sync.WaitGroup)

				for _, component := range c.application.GetComponents() {
					if capture, ok := component.(hasCapture); ok {
						wg.Add(1)

						go func() {
							defer wg.Done()
							capture.MetricsCapture()
						}()
					}
				}

				wg.Wait()

				// send metrics
				c.mutex.RLock()
				client := c.client
				c.mutex.RUnlock()

				if err := c.connector.WriteTo(client); err != nil {
					c.logger.Error("Send metric to Influx failed", map[string]interface{}{
						"error": err.Error(),
					})
				} else {
					c.logger.Debug("Send metric to Influx success")
				}
			case d := <-c.changeTicker:
				ticker = time.NewTicker(d)
			}
		}
	}()

	return nil
}

func (c *Component) initClient(url, username, password string) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.client, err = influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     url,
		Username: username,
		Password: password,
	})

	return err
}

func (c *Component) getName(name string) string {
	return fmt.Sprint(c.prefix, name)
}

func (c *Component) NewCounter(name string) kit.Counter {
	return c.connector.NewCounter(c.getName(name))
}

func (c *Component) NewGauge(name string) kit.Gauge {
	return c.connector.NewGauge(c.getName(name))
}

func (c *Component) NewHistogram(name string) kit.Histogram {
	return c.connector.NewHistogram(c.getName(name))
}

func (c *Component) NewTimer(name string) Timer {
	return NewMetricTimer(c.NewHistogram(name))
}

func (c *Component) CaptureMetrics(d time.Duration, f func()) {
	go func() {
		for range time.NewTicker(d).C {
			f()
		}
	}()
}

func (c *Component) getTags() map[string]string {
	tags := map[string]string{
		TagAppName:    c.application.GetName(),
		TagAppVersion: c.application.GetVersion(),
		TagAppBuild:   c.application.GetBuild(),
	}

	if hostname, err := os.Hostname(); err == nil {
		tags[TagHostname] = hostname
	}

	tagsFromConfig := c.config.GetString(ConfigMetricsTags)
	if len(tagsFromConfig) > 0 {
		var parts []string

		for _, tag := range strings.Split(tagsFromConfig, ",") {
			parts = strings.Split(tag, "=")

			if len(parts) > 1 {
				tags[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	return tags
}
