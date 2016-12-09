package frontend

import (
	"github.com/kihamo/shadow/resource/metrics"
)

const (
	MetricHandlerExecuteTime = "handler_execute_time"
)

var (
	metricHandlerExecuteTime metrics.Timer
)

func (s *FrontendService) MetricsRegister(m *metrics.Resource) {
	metricHandlerExecuteTime = m.NewTimer(MetricHandlerExecuteTime)
}
