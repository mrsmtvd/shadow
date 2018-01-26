package internal

import (
	"github.com/kihamo/shadow/components/database"
	"github.com/kihamo/shadow/components/database/storage"
	"github.com/kihamo/snitch"
)

const (
	MetricOpenConnectionsTotal = database.ComponentName + "_open_connections_total"
)

var (
	metricOpenConnectionsTotal snitch.Gauge
)

type metricsCollector struct {
	component *Component
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	metricOpenConnectionsTotal.Describe(ch)

	// describe from storages
	storage.Describe(ch)
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	s := c.component.Storage()

	if s == nil {
		return
	}

	stats := s.Master().(*storage.SQLExecutor).DB().Stats()

	metricOpenConnectionsTotal.Set(float64(stats.OpenConnections))

	metricOpenConnectionsTotal.Collect(ch)

	// collect from storages
	storage.CollectStorageSQL(ch)
}

func (c *Component) Metrics() snitch.Collector {
	metricOpenConnectionsTotal = snitch.NewGauge(MetricOpenConnectionsTotal, "Number of open connections to the database")

	return &metricsCollector{
		component: c,
	}
}
