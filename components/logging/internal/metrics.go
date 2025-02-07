package internal

import (
	"github.com/kihamo/snitch"
	"github.com/mrsmtvd/shadow/components/logging"
)

const (
	MetricTotal = logging.ComponentName + "_send_total"
)

var (
	metricTotal = snitch.NewCounter(MetricTotal, "Number of send logs")
)

func (c *Component) Metrics() snitch.Collector {
	return metricTotal
}
