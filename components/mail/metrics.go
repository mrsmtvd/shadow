package mail

import (
	"github.com/kihamo/snitch"
)

const (
	MetricTotal = ComponentName + "_send_total"
)

var (
	metricMailTotal snitch.Counter
)

func (c *Component) Metrics() snitch.Collector {
	metricMailTotal = snitch.NewCounter(MetricTotal, "Number of mail")

	return metricMailTotal
}
