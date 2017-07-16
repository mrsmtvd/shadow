package grpc

import (
	"github.com/kihamo/snitch"
)

const (
	MetricExecuteTime       = ComponentName + "_request_duration_seconds"
	MetricMethodExecuteTime = ComponentName + "_method_duration_seconds"
)

var (
	metricExecuteTime snitch.Timer
)

func (c *Component) Metrics() snitch.Collector {
	metricExecuteTime = snitch.NewTimer(MetricExecuteTime, "Total request duration")

	return metricExecuteTime
}
