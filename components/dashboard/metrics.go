package dashboard

import (
	"github.com/kihamo/snitch"
)

const (
	MetricHandlerExecuteTime = ComponentName + "_handler_response_time_milliseconds"
)

var (
	metricHandlerExecuteTime = snitch.NewTimer(MetricHandlerExecuteTime, "Time of handle requests")
)

func (c *Component) Metrics() snitch.Collector {
	return metricHandlerExecuteTime
}
