package internal

import (
	"github.com/kihamo/shadow/components/grpc"
	"github.com/kihamo/snitch"
)

const (
	MetricExecuteTime = grpc.ComponentName + "_request_duration_seconds"
)

var (
	metricExecuteTime snitch.Timer
)

func (c *Component) Metrics() snitch.Collector {
	metricExecuteTime = snitch.NewTimer(MetricExecuteTime, "Total request duration")

	return metricExecuteTime
}
