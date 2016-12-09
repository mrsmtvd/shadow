package mail

import (
	kit "github.com/go-kit/kit/metrics"
	"github.com/kihamo/shadow/resource/metrics"
)

const (
	MetricMailTotal = "mail.total"
)

var (
	metricMailTotal kit.Counter
)

func (r *Resource) MetricsRegister(m *metrics.Resource) {
	metricMailTotal = m.NewCounter(MetricMailTotal)
}
