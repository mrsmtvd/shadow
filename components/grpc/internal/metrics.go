package internal

import (
	"github.com/kihamo/shadow/components/grpc/stats"
	"github.com/kihamo/snitch"
)

type metricsCollector struct {
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	stats.MetricHandledTotal.Describe(ch)
	stats.MetricReceivedTotal.Describe(ch)
	stats.MetricSentTotal.Describe(ch)
	stats.MetricStartedTotal.Describe(ch)
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	stats.MetricHandledTotal.Collect(ch)
	stats.MetricReceivedTotal.Collect(ch)
	stats.MetricSentTotal.Collect(ch)
	stats.MetricStartedTotal.Collect(ch)
}

func (c *Component) Metrics() snitch.Collector {
	return &metricsCollector{}
}
