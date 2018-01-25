package internal

import (
	"time"

	"github.com/kihamo/shadow/components/database"
	"github.com/kihamo/shadow/components/database/storage"
	"github.com/kihamo/snitch"
)

const (
	MetricOpenConnectionsTotal = database.ComponentName + "_open_connections_total"
	MetricQueryDuration        = database.ComponentName + "_query_duration_seconds"

	OperationExec   = "exec"
	OperationCreate = "create"
	OperationSelect = "select"
	OperationInsert = "insert"
	OperationUpdate = "update"
	OperationDelete = "delete"
)

var (
	metricOpenConnectionsTotal snitch.Gauge
	metricQueryDuration        snitch.Timer
)

type metricsCollector struct {
	component *Component
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	metricOpenConnectionsTotal.Describe(ch)
	metricQueryDuration.Describe(ch)
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	s := c.component.Storage()

	if s == nil {
		return
	}

	stats := s.Master().(*storage.SQLExecutor).DB().Stats()

	metricOpenConnectionsTotal.Set(float64(stats.OpenConnections))

	metricOpenConnectionsTotal.Collect(ch)
	metricQueryDuration.Collect(ch)
}

func (c *Component) Metrics() snitch.Collector {
	metricOpenConnectionsTotal = snitch.NewGauge(MetricOpenConnectionsTotal, "Number of open connections to the database")
	metricQueryDuration = snitch.NewTimer(MetricQueryDuration, "Response time of queries to the database")

	return &metricsCollector{
		component: c,
	}
}

func updateMetric(operation string, startAt time.Time) {
	if metricQueryDuration != nil {
		metricQueryDuration.With("type", operation).UpdateSince(startAt)
	}
}
