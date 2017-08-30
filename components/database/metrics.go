package database

import (
	"github.com/go-gorp/gorp"
	"github.com/kihamo/snitch"
)

const (
	MetricOpenConnectionsTotal = ComponentName + "_open_connections_total"
)

var (
	metricOpenConnectionsTotal snitch.Gauge
)

type metricsCollector struct {
	component *Component
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	ch <- metricOpenConnectionsTotal.Description()
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	storage := c.component.GetStorage()

	if storage == nil {
		return
	}

	stats := storage.executor.(*gorp.DbMap).Db.Stats()

	metricOpenConnectionsTotal.Set(float64(stats.OpenConnections))

	ch <- metricOpenConnectionsTotal
}

func (c *Component) Metrics() snitch.Collector {
	metricOpenConnectionsTotal = snitch.NewGauge(MetricOpenConnectionsTotal, "Number of open connections to the database")

	return &metricsCollector{
		component: c,
	}
}
