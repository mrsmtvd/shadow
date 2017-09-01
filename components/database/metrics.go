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

	metricQueriesTotal       snitch.Counter
	metricQueriesTotalExec   snitch.Counter
	metricQueriesTotalCreate snitch.Counter
	metricQueriesTotalSelect snitch.Counter
	metricQueriesTotalInsert snitch.Counter
	metricQueriesTotalUpdate snitch.Counter
	metricQueriesTotalDelete snitch.Counter

	metricQueryDuration       snitch.Timer
	metricQueryDurationExec   snitch.Timer
	metricQueryDurationCreate snitch.Timer
	metricQueryDurationSelect snitch.Timer
	metricQueryDurationInsert snitch.Timer
	metricQueryDurationUpdate snitch.Timer
	metricQueryDurationDelete snitch.Timer
)

type metricsCollector struct {
	component *Component
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	ch <- metricOpenConnectionsTotal.Description()

	ch <- metricQueriesTotal.Description()
	ch <- metricQueriesTotalExec.Description()
	ch <- metricQueriesTotalCreate.Description()
	ch <- metricQueriesTotalSelect.Description()
	ch <- metricQueriesTotalInsert.Description()
	ch <- metricQueriesTotalUpdate.Description()
	ch <- metricQueriesTotalDelete.Description()

	ch <- metricQueryDuration.Description()
	ch <- metricQueryDurationExec.Description()
	ch <- metricQueryDurationCreate.Description()
	ch <- metricQueryDurationSelect.Description()
	ch <- metricQueryDurationInsert.Description()
	ch <- metricQueryDurationUpdate.Description()
	ch <- metricQueryDurationDelete.Description()
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	storage := c.component.GetStorage()

	if storage == nil {
		return
	}

	stats := storage.executor.(*gorp.DbMap).Db.Stats()

	metricOpenConnectionsTotal.Set(float64(stats.OpenConnections))

	ch <- metricOpenConnectionsTotal

	ch <- metricQueriesTotal
	ch <- metricQueriesTotalExec
	ch <- metricQueriesTotalCreate
	ch <- metricQueriesTotalSelect
	ch <- metricQueriesTotalInsert
	ch <- metricQueriesTotalUpdate
	ch <- metricQueriesTotalDelete

	ch <- metricQueryDuration
	ch <- metricQueryDurationExec
	ch <- metricQueryDurationCreate
	ch <- metricQueryDurationSelect
	ch <- metricQueryDurationInsert
	ch <- metricQueryDurationUpdate
	ch <- metricQueryDurationDelete
}

func (c *Component) Metrics() snitch.Collector {
	metricOpenConnectionsTotal = snitch.NewGauge(MetricOpenConnectionsTotal, "Number of open connections to the database")

	metricQueriesTotal = snitch.NewCounter(MetricQueriesTotal, "Number of queries to the database")
	metricQueriesTotalExec = snitch.NewCounter(MetricQueriesTotal, "Number of exec queries to the database", "type", OperationExec)
	metricQueriesTotalCreate = snitch.NewCounter(MetricQueriesTotal, "Number of select queries to the database", "type", OperationCreate)
	metricQueriesTotalSelect = snitch.NewCounter(MetricQueriesTotal, "Number of select queries to the database", "type", OperationSelect)
	metricQueriesTotalInsert = snitch.NewCounter(MetricQueriesTotal, "Number of insert queries to the database", "type", OperationInsert)
	metricQueriesTotalUpdate = snitch.NewCounter(MetricQueriesTotal, "Number of update queries to the database", "type", OperationUpdate)
	metricQueriesTotalDelete = snitch.NewCounter(MetricQueriesTotal, "Number of delete queries to the database", "type", OperationDelete)

	metricQueryDuration = snitch.NewTimer(MetricQueryDuration, "Response time of queries to the database")
	metricQueryDurationExec = snitch.NewTimer(MetricQueryDuration, "Response time of exec queries to the database", "type", OperationExec)
	metricQueryDurationCreate = snitch.NewTimer(MetricQueryDuration, "Response time of select queries to the database", "type", OperationCreate)
	metricQueryDurationSelect = snitch.NewTimer(MetricQueryDuration, "Response time of select queries to the database", "type", OperationSelect)
	metricQueryDurationInsert = snitch.NewTimer(MetricQueryDuration, "Response time of insert queries to the database", "type", OperationInsert)
	metricQueryDurationUpdate = snitch.NewTimer(MetricQueryDuration, "Response time of update queries to the database", "type", OperationUpdate)
	metricQueryDurationDelete = snitch.NewTimer(MetricQueryDuration, "Response time of delete queries to the database", "type", OperationDelete)

	return &metricsCollector{
		component: c,
	}
}

func updateMetric(operation string, startAt time.Time) {
	if metricQueriesTotal != nil {
		metricQueriesTotal.Inc()
	}

	if metricQueryDuration != nil {
		metricQueryDuration.UpdateSince(startAt)
	}

	var (
		total    snitch.Counter
		duration snitch.Timer
	)

	switch operation {
	case OperationCreate:
		total = metricQueriesTotalCreate
		duration = metricQueryDurationCreate
	case OperationSelect:
		total = metricQueriesTotalSelect
		duration = metricQueryDurationSelect
	case OperationInsert:
		total = metricQueriesTotalInsert
		duration = metricQueryDurationInsert
	case OperationUpdate:
		total = metricQueriesTotalUpdate
		duration = metricQueryDurationUpdate
	case OperationDelete:
		total = metricQueriesTotalDelete
		duration = metricQueryDurationDelete
	default:
		total = metricQueriesTotalExec
		duration = metricQueryDurationExec
	}

	if total != nil {
		total.Inc()
	}

	if duration != nil {
		duration.UpdateSince(startAt)
	}
}
