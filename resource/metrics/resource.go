package metrics

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/influx"
	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
	"github.com/kihamo/shadow/resource/logger"
)

type Resource struct {
	application *shadow.Application

	config *config.Resource
	logger *logger.Resource

	connector *influx.Influx
	prefix    string
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

	var l log.Logger

	if r.application.HasResource("logger") {
		resourceLogger, _ := r.application.GetResource("logger")
		l = newMetricsLogger(resourceLogger.(*logger.Resource).Get("metrics"))
	} else {
		l = log.NewNopLogger()
	}

	r.connector = influx.New(r.getTags(), influxdb.BatchPointsConfig{
		Database: r.config.GetString("metrics.database"),
		Precision: "s",
	}, l)

	r.prefix = r.config.GetString("metrics.prefix")

	go func() {
		defer wg.Done()

		ticker := time.NewTicker(r.config.GetDuration("metrics.interval"))
		defer ticker.Stop()

		r.connector.WriteLoop(ticker.C, client)
	}()

	return nil
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
