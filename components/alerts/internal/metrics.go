package internal

import (
	"github.com/kihamo/shadow/components/alerts"
	"github.com/kihamo/snitch"
)

const (
	MetricTotal = alerts.ComponentName + "_send_total"
)

var (
	metricAlertsTotal snitch.Counter
)

func (c *Component) Metrics() snitch.Collector {
	metricAlertsTotal = snitch.NewCounter(MetricTotal, "Number of alerts")

	return metricAlertsTotal
}
