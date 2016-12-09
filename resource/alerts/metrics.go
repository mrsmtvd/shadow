package alerts

import (
	kit "github.com/go-kit/kit/metrics"
	"github.com/kihamo/shadow/resource/metrics"
)

const (
	MetricAlertsTotal = "alerts.total"
)

var (
	metricAlertsTotal kit.Counter
)

func (r *Resource) MetricsRegister(m *metrics.Resource) {
	metricAlertsTotal = m.NewCounter(MetricAlertsTotal)
}
