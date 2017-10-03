package internal

import (
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/snitch"
)

const (
	MetricHandlerExecuteTime = dashboard.ComponentName + "_handler_response_time_seconds"
)

var (
	metricHandlerExecuteTime = snitch.NewTimer(MetricHandlerExecuteTime, "Time of handle requests")
)

func (c *Component) Metrics() snitch.Collector {
	return metricHandlerExecuteTime
}
