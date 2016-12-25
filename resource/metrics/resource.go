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
	"github.com/kihamo/shadow/resource/config"
	"github.com/kihamo/shadow/resource/logger"
)

const (
	TagAppName    = "app_name"
	TagAppVersion = "app_version"
	TagAppBuild   = "app_build"
	TagHostname   = "hostname"
)

type Resource struct {
	application *shadow.Application

	config *config.Resource
	logger logger.Logger

	mutex        sync.RWMutex
	client       influxdb.Client
	connector    *influx.Influx
	prefix       string
	changeTicker chan time.Duration
}

type ContextItemMetrics interface {
	MetricsRegister(r *Resource)
}

type ContextItemMetricsCapture interface {
	MetricsCapture()
}

func (r *Resource) GetName() string {
	return "metrics"
}

func (r *Resource) Init(a *shadow.Application) error {
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	r.config = resourceConfig.(*config.Resource)

	r.application = a

	r.changeTicker = make(chan time.Duration)

	return nil
}

func (r *Resource) Run(wg *sync.WaitGroup) error {
	if err := r.initClient(r.config.GetString(ConfigMetricsUrl), r.config.GetString(ConfigMetricsUsername), r.config.GetString(ConfigMetricsPassword)); err != nil {
		return err
	}

	r.logger = logger.NewOrNop(r.GetName(), r.application)

	r.connector = influx.New(r.getTags(), influxdb.BatchPointsConfig{
		Database:  r.config.GetString(ConfigMetricsDatabase),
		Precision: r.config.GetString(ConfigMetricsPrecision),
	}, logger.NewGoKitLogger(r.logger))

	r.prefix = r.config.GetString(ConfigMetricsPrefix)

	// search metrics
	for _, resource := range r.application.GetResources() {
		if rMetrics, ok := resource.(ContextItemMetrics); ok {
			rMetrics.MetricsRegister(r)
		}
	}

	for _, service := range r.application.GetServices() {
		if sMetrics, ok := service.(ContextItemMetrics); ok {
			sMetrics.MetricsRegister(r)
		}
	}

	// send to influx
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(r.config.GetDuration(ConfigMetricsInterval))

		for {
			select {
			case <-ticker.C:
				// collect metrics
				wg := new(sync.WaitGroup)

				for _, resource := range r.application.GetResources() {
					if rCapture, ok := resource.(ContextItemMetricsCapture); ok {
						wg.Add(1)

						go func() {
							defer wg.Done()
							rCapture.MetricsCapture()
						}()
					}
				}

				for _, service := range r.application.GetServices() {
					if sCapture, ok := service.(ContextItemMetricsCapture); ok {
						wg.Add(1)

						go func() {
							defer wg.Done()
							sCapture.MetricsCapture()
						}()
					}
				}

				wg.Wait()

				// send metrics
				r.mutex.RLock()
				client := r.client
				r.mutex.RUnlock()

				if err := r.connector.WriteTo(client); err != nil {
					r.logger.Error("Send metric to Influx failed", map[string]interface{}{
						"error": err.Error(),
					})
				} else {
					r.logger.Debug("Send metric to Influx success")
				}
			case d := <-r.changeTicker:
				ticker = time.NewTicker(d)
			}
		}
	}()

	return nil
}

func (r *Resource) initClient(url, username, password string) (err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.client, err = influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     url,
		Username: username,
		Password: password,
	})

	return err
}

func (r *Resource) getName(name string) string {
	return fmt.Sprint(r.prefix, name)
}

func (r *Resource) NewCounter(name string) kit.Counter {
	return r.connector.NewCounter(r.getName(name))
}

func (r *Resource) NewGauge(name string) kit.Gauge {
	return r.connector.NewGauge(r.getName(name))
}

func (r *Resource) NewHistogram(name string) kit.Histogram {
	return r.connector.NewHistogram(r.getName(name))
}

func (r *Resource) NewTimer(name string) Timer {
	return NewMetricTimer(r.NewHistogram(name))
}

func (r *Resource) CaptureMetrics(d time.Duration, f func()) {
	go func() {
		for range time.NewTicker(d).C {
			f()
		}
	}()
}

func (r *Resource) getTags() map[string]string {
	tags := map[string]string{
		TagAppName:    r.application.Name,
		TagAppVersion: r.application.Version,
		TagAppBuild:   r.application.Build,
	}

	if hostname, err := os.Hostname(); err == nil {
		tags[TagHostname] = hostname
	}

	tagsFromConfig := r.config.GetString(ConfigMetricsTags)
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
