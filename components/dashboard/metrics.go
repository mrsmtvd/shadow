package dashboard

import (
	"github.com/kihamo/snitch"
)

const (
	MetricHandlerExecuteTime = ComponentName + ".handler.execute_time"
)

var (
	metricHandlerExecuteTime snitch.Timer
)

func (c *Component) Metrics() snitch.Collector {
	metricHandlerExecuteTime = snitch.NewTimer(MetricHandlerExecuteTime)

	return metricHandlerExecuteTime
}
