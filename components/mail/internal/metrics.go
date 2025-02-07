package internal

import (
	"github.com/kihamo/snitch"
	"github.com/mrsmtvd/shadow/components/mail"
)

const (
	MetricTotal = mail.ComponentName + "_send_total"
)

var (
	metricMailTotal = snitch.NewCounter(MetricTotal, "Number of mail")
)

func (c *Component) Metrics() snitch.Collector {
	return metricMailTotal
}
