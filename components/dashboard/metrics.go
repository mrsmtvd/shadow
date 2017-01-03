package dashboard

import (
	"github.com/kihamo/shadow/components/metrics"
)

const (
	MetricHandlerExecuteTime = "handler_execute_time"
)

var (
	metricHandlerExecuteTime metrics.Timer
)

func (c *Component) MetricsRegister(m *metrics.Component) {
	metricHandlerExecuteTime = m.NewTimer(MetricHandlerExecuteTime)
}
