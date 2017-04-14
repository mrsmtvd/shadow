package alerts

import (
	"github.com/kihamo/snitch"
)

const (
	MetricAlertsTotal = "alerts.total"
)

var (
	metricAlertsTotal snitch.Counter
)

func (c *Component) Metrics() snitch.Collector {
	metricAlertsTotal = snitch.NewCounter(MetricAlertsTotal)

	return metricAlertsTotal
}
