package internal

import (
	"github.com/kihamo/snitch"
	"github.com/mrsmtvd/shadow/components/dashboard"
)

const (
	MetricHealthCheckStatus = dashboard.ComponentName + "_healthcheck_status"
)

var (
	metricHealthCheckStatus = snitch.NewGauge(MetricHealthCheckStatus, "Current check status (0 indicates success, 1 indicates failure)")
)

func (c *Component) Metrics() snitch.Collector {
	return metricHealthCheckStatus
}
