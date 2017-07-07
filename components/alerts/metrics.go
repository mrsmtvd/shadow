package alerts

import (
	"github.com/kihamo/snitch"
)

const (
	MetricTotal = ComponentName + "_send_total"
)

var (
	metricAlertsTotal snitch.Counter
)

func (c *Component) Metrics() snitch.Collector {
	metricAlertsTotal = snitch.NewCounter(MetricTotal, "Number of alerts")

	return metricAlertsTotal
}
