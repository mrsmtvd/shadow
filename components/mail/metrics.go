package mail

import (
	"github.com/kihamo/snitch"
)

const (
	MetricTotal = ComponentName + "_send_total"
)

var (
	metricMailTotalSuccess snitch.Counter
	metricMailTotalFailed  snitch.Counter
)

type metricsCollector struct {
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	ch <- metricMailTotalSuccess.Description()
	ch <- metricMailTotalFailed.Description()
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	ch <- metricMailTotalSuccess
	ch <- metricMailTotalFailed
}

func (c *Component) Metrics() snitch.Collector {
	metricMailTotalSuccess = snitch.NewCounter(MetricTotal, "status", "success")
	metricMailTotalFailed = snitch.NewCounter(MetricTotal, "status", "failed")

	return &metricsCollector{}
}
