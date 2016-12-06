package metrics

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
	"github.com/kihamo/shadow/resource/logger"
	"github.com/rcrowley/go-metrics"
	"github.com/vrischmann/go-metrics-influxdb"
)

type Metrics struct {
	application *shadow.Application
	config      *config.Config
	registry    metrics.Registry
}

func (r *Metrics) GetName() string {
	return "metrics"
}

func (r *Metrics) GetConfigVariables() []config.ConfigVariable {
	return []config.ConfigVariable{
		config.ConfigVariable{
			Key:   "metrics.url",
			Value: "",
			Usage: "InfluxDB url",
		},
		config.ConfigVariable{
			Key:   "metrics.database",
			Value: "metrics",
			Usage: "InfluxDB database name",
		},
		config.ConfigVariable{
			Key:   "metrics.username",
			Value: "",
			Usage: "InfluxDB username",
		},
		config.ConfigVariable{
			Key:   "metrics.password",
			Value: "",
			Usage: "InfluxDB password",
		},
		config.ConfigVariable{
			Key:   "metrics.interval",
			Value: "20s",
			Usage: "Flush interval",
		},
		config.ConfigVariable{
			Key:   "metrics.tags",
			Value: "",
			Usage: "Tags list with format: tag1_name=tag1_value,tag2_name=tag2_value",
		},
	}
}

func (r *Metrics) Init(a *shadow.Application) error {
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	r.config = resourceConfig.(*config.Config)

	r.application = a

	return nil
}

func (r *Metrics) Run(wg *sync.WaitGroup) error {
	registry := r.getRegistry()

	if r.config.GetBool("debug") {
		metrics.RegisterDebugGCStats(registry)
		go func() {
			defer wg.Done()
			metrics.CaptureDebugGCStats(registry, r.config.GetDuration("metrics.interval"))
		}()
	}

	metrics.RegisterRuntimeMemStats(registry)
	go func() {
		defer wg.Done()
		metrics.CaptureRuntimeMemStats(registry, r.config.GetDuration("metrics.interval"))
	}()

	return nil
}

func (r *Metrics) getRegistry() metrics.Registry {
	if r.registry != nil {
		return r.registry
	}

	r.registry = metrics.NewRegistry()

	go func() {
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

		influxdb.InfluxDBWithTags(
			r.registry,
			r.config.GetDuration("metrics.interval"),
			r.config.GetString("metrics.url"),
			r.config.GetString("metrics.database"),
			r.config.GetString("metrics.username"),
			r.config.GetString("metrics.password"),
			tags)

	}()

	if r.config.GetBool("debug") {
		resourceLogger, err := r.application.GetResource("logger")
		if err == nil {
			go func() {
				resourceLogger := resourceLogger.(*logger.Logger).Get(r.GetName())
				metricsLogger := log.New(resourceLogger, "", 0)

				metrics.Log(
					r.registry,
					r.config.GetDuration("metrics.interval"),
					metricsLogger)
			}()
		}
	}

	return r.registry
}
