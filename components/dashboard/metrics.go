package dashboard

import (
	"github.com/kihamo/snitch"
)

const (
	MetricHandlerExecuteTime = ComponentName + "_handler_response_time_milliseconds"
)

var (
	metricHandlerExecuteTime snitch.Timer
)

func (c *Component) Metrics() snitch.Collector {
	metricHandlerExecuteTime = snitch.NewTimer(MetricHandlerExecuteTime)

	return metricHandlerExecuteTime
}
