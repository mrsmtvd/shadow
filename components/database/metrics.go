package database

import (
	"time"

	"github.com/go-gorp/gorp"
	"github.com/kihamo/snitch"
)

const (
	MetricOpenConnectionsTotal = ComponentName + "_open_connections_total"
	MetricQueriesTotal         = ComponentName + "_queries_total"
	MetricQueryDuration        = ComponentName + "_query_duration_seconds"

	OperationExec   = "exec"
	OperationCreate = "create"
	OperationSelect = "select"
	OperationInsert = "insert"
	OperationUpdate = "update"
	OperationDelete = "delete"
)

var (
	metricOpenConnectionsTotal snitch.Gauge
	metricQueriesTotal         snitch.Counter
	metricQueryDuration        snitch.Timer
)

type metricsCollector struct {
	component *Component
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	metricOpenConnectionsTotal.Describe(ch)
	metricQueriesTotal.Describe(ch)
	metricQueryDuration.Describe(ch)
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	storage := c.component.GetStorage()

	if storage == nil {
		return
	}

	stats := storage.executor.(*gorp.DbMap).Db.Stats()

	metricOpenConnectionsTotal.Set(float64(stats.OpenConnections))

	metricOpenConnectionsTotal.Collect(ch)
	metricQueriesTotal.Collect(ch)
	metricQueryDuration.Collect(ch)
}

func (c *Component) Metrics() snitch.Collector {
	metricOpenConnectionsTotal = snitch.NewGauge(MetricOpenConnectionsTotal, "Number of open connections to the database")
	metricQueriesTotal = snitch.NewCounter(MetricQueriesTotal, "Number of queries to the database")
	metricQueryDuration = snitch.NewTimer(MetricQueryDuration, "Response time of queries to the database")

	return &metricsCollector{
		component: c,
	}
}

func updateMetric(operation string, startAt time.Time) {
	if metricQueriesTotal != nil {
		metricQueriesTotal.Inc()
		metricQueriesTotal.With("type", operation).Inc()
	}

	if metricQueryDuration != nil {
		metricQueryDuration.UpdateSince(startAt)
		metricQueryDuration.With("type", operation).UpdateSince(startAt)
	}
}
