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

type Resource struct {
	application *shadow.Application

	config *config.Resource
	logger logger.Logger

	connector *influx.Influx
	prefix    string
}

type ContextItemMetrics interface {
	MetricsRegister(r *Resource)
}

type ContextItemMetricsCapture interface {
	MetricsCapture() (func(), time.Duration)
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

	return nil
}

func (r *Resource) Run(wg *sync.WaitGroup) error {
	client, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     r.config.GetString("metrics.url"),
		Username: r.config.GetString("metrics.username"),
		Password: r.config.GetString("metrics.password"),
	})

	if err != nil {
		return err
	}

	if resourceLogger, err := r.application.GetResource("logger"); err == nil {
		r.logger = resourceLogger.(*logger.Resource).Get(r.GetName())
	} else {
		r.logger = logger.NopLogger
	}

	r.connector = influx.New(r.getTags(), influxdb.BatchPointsConfig{
		Database:  r.config.GetString("metrics.database"),
		Precision: "s",
	}, r.logger)

	r.prefix = r.config.GetString("metrics.prefix")

	// send to influx
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(r.config.GetDuration("metrics.interval"))
		defer ticker.Stop()

		for range ticker.C {
			if err := r.connector.WriteTo(client); err != nil {
				r.logger.Error("Send metric to Influx failed")
			} else {
				r.logger.Debug("Send metric to Influx success")
			}
		}
	}()

	// search metrics
	for _, resource := range r.application.GetResources() {
		if rMetrics, ok := resource.(ContextItemMetrics); ok {
			rMetrics.MetricsRegister(r)
		}

		if rCapture, ok := resource.(ContextItemMetricsCapture); ok {
			f, d := rCapture.MetricsCapture()
			r.CaptureMetrics(d, f)
		}
	}

	for _, service := range r.application.GetServices() {
		if sMetrics, ok := service.(ContextItemMetrics); ok {
			sMetrics.MetricsRegister(r)
		}

		if sCapture, ok := service.(ContextItemMetricsCapture); ok {
			f, d := sCapture.MetricsCapture()
			r.CaptureMetrics(d, f)
		}
	}

	return nil
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
		"app_name":    r.application.Name,
		"app_version": r.application.Version,
		"app_build":   r.application.Build,
	}

	if hostname, err := os.Hostname(); err == nil {
		tags["hostname"] = hostname
	}

	tagsFromConfig := r.config.GetString("metrics.tags")
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
