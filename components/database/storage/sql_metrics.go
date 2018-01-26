package storage

import (
	"time"

	"github.com/kihamo/shadow/components/database"
	"github.com/kihamo/snitch"
)

const (
	MetricQueryDuration = database.ComponentName + "_query_duration_seconds"

	OperationExec   = "exec"
	OperationCreate = "create"
	OperationSelect = "select"
	OperationInsert = "insert"
	OperationUpdate = "update"
	OperationDelete = "delete"
)

var metricStorageSQLQueryDuration = snitch.NewTimer(MetricQueryDuration, "Response time of queries to the database")

func Describe(ch chan<- *snitch.Description) {
	metricStorageSQLQueryDuration.Describe(ch)
}

func CollectStorageSQL(ch chan<- snitch.Metric) {
	metricStorageSQLQueryDuration.Collect(ch)
}

func UpdateStorageSQLMetric(operation, server string, startAt time.Time) {
	metricStorageSQLQueryDuration.With(
		"operation", operation,
		"server", server,
	).UpdateSince(startAt)
}
