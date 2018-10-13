package internal

import (
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/snitch"
)

type metricsCollector struct {
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	metrics.MetricResponseTimeSeconds.Describe(ch)
	metrics.MetricResponseSizeBytes.Describe(ch)
	metrics.MetricResponseMarshalTimeSeconds.Describe(ch)

	metrics.MetricRequestSizeBytes.Describe(ch)
	metrics.MetricRequestsTotal.Describe(ch)
	metrics.MetricRequestReadTimeSeconds.Describe(ch)
	metrics.MetricRequestReadUnmarshalTimeSeconds.Describe(ch)
	metrics.MetricRequestUnmarshalTimeSeconds.Describe(ch)

	metrics.MetricExternalResponseTimeSeconds.Describe(ch)

	metrics.MetricGRPCHandledTotal.Describe(ch)
	metrics.MetricGRPCReceivedTotal.Describe(ch)
	metrics.MetricGRPCSentTotal.Describe(ch)
	metrics.MetricGRPCStartedTotal.Describe(ch)
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	metrics.MetricResponseTimeSeconds.Collect(ch)
	metrics.MetricResponseSizeBytes.Collect(ch)
	metrics.MetricResponseMarshalTimeSeconds.Collect(ch)

	metrics.MetricRequestSizeBytes.Collect(ch)
	metrics.MetricRequestsTotal.Collect(ch)
	metrics.MetricRequestReadTimeSeconds.Collect(ch)
	metrics.MetricRequestReadUnmarshalTimeSeconds.Collect(ch)
	metrics.MetricRequestUnmarshalTimeSeconds.Collect(ch)

	metrics.MetricExternalResponseTimeSeconds.Collect(ch)

	metrics.MetricGRPCHandledTotal.Collect(ch)
	metrics.MetricGRPCReceivedTotal.Collect(ch)
	metrics.MetricGRPCSentTotal.Collect(ch)
	metrics.MetricGRPCStartedTotal.Collect(ch)
}

func (c *Component) Metrics() snitch.Collector {
	return &metricsCollector{}
}
