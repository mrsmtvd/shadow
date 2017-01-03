package alerts

import (
	kit "github.com/go-kit/kit/metrics"
	"github.com/kihamo/shadow/components/metrics"
)

const (
	MetricAlertsTotal = "alerts.total"
)

var (
	metricAlertsTotal kit.Counter
)

func (c *Component) MetricsRegister(m *metrics.Component) {
	metricAlertsTotal = m.NewCounter(MetricAlertsTotal)
}
