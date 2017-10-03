package internal

import (
	"github.com/kihamo/shadow/components/mail"
	"github.com/kihamo/snitch"
)

const (
	MetricTotal = mail.ComponentName + "_send_total"
)

var (
	metricMailTotal snitch.Counter
)

func (c *Component) Metrics() snitch.Collector {
	metricMailTotal = snitch.NewCounter(MetricTotal, "Number of mail")

	return metricMailTotal
}
