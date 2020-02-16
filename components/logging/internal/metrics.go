package internal

import (
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/snitch"
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
