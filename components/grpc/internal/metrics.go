package internal

import (
	"github.com/kihamo/shadow/components/grpc"
	"github.com/kihamo/snitch"
)

const (
	MetricRequestDuration = grpc.ComponentName + "_request_duration_seconds"
	MetricRequests        = grpc.ComponentName + "_requests_total"
)

var (
	metricRequestDuration = snitch.NewTimer(MetricRequestDuration, "Total request duration")
	metricRequests        = snitch.NewCounter(MetricRequests, "Total requests")
)

type metricsCollector struct {
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	metricRequestDuration.Describe(ch)
	metricRequests.Describe(ch)
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	metricRequestDuration.Collect(ch)
	metricRequests.Collect(ch)
}

func (c *Component) Metrics() snitch.Collector {
	return &metricsCollector{}
}
