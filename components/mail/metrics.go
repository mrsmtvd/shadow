package mail

import (
	kit "github.com/go-kit/kit/metrics"
	"github.com/kihamo/shadow/components/metrics"
)

const (
	MetricMailTotal = "mail.total"
)

var (
	metricMailTotal kit.Counter
)

func (c *Component) MetricsRegister(m *metrics.Component) {
	metricMailTotal = m.NewCounter(MetricMailTotal)
}
