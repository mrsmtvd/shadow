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

	stats.MetricResponseTimeSeconds.Describe(ch)
	stats.MetricRequestsTotal.Describe(ch)
	stats.MetricRequestSizeBytes.Describe(ch)
	stats.MetricResponseSizeBytes.Describe(ch)
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	stats.MetricHandledTotal.Collect(ch)
	stats.MetricReceivedTotal.Collect(ch)
	stats.MetricSentTotal.Collect(ch)
	stats.MetricStartedTotal.Collect(ch)

	stats.MetricResponseTimeSeconds.Collect(ch)
	stats.MetricRequestsTotal.Collect(ch)
	stats.MetricRequestSizeBytes.Collect(ch)
	stats.MetricResponseSizeBytes.Collect(ch)
}

func (c *Component) Metrics() snitch.Collector {
	return &metricsCollector{}
}
