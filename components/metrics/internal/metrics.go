package internal

import (
	"time"

	"github.com/kihamo/snitch"
	"github.com/mrsmtvd/shadow"
	"github.com/mrsmtvd/shadow/components/metrics"
)

type metricsCollector struct {
	application shadow.Application

	upTime    snitch.Counter
	buildInfo snitch.Gauge
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	metrics.MetricResponseTimeSeconds.Describe(ch)
	metrics.MetricResponseSizeBytes.Describe(ch)
	metrics.MetricResponseMarshalTimeSeconds.Describe(ch)

	metrics.MetricRequestSizeBytes.Describe(ch)
	metrics.MetricRequestsTotal.Describe(ch)
	metrics.MetricRequestReadTimeSeconds.Describe(ch)
	metrics.MetricRequestReadUnmarshalTimeSeconds.Describe(ch)
	metrics.MetricRequestUnmarshalTimeSeconds.Describe(ch)

	metrics.MetricExternalResponseTimeSeconds.Describe(ch)

	metrics.MetricGRPCHandledTotal.Describe(ch)
	metrics.MetricGRPCReceivedTotal.Describe(ch)
	metrics.MetricGRPCSentTotal.Describe(ch)
	metrics.MetricGRPCStartedTotal.Describe(ch)

	c.upTime.Describe(ch)
	c.buildInfo.Describe(ch)
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	metrics.MetricResponseTimeSeconds.Collect(ch)
	metrics.MetricResponseSizeBytes.Collect(ch)
	metrics.MetricResponseMarshalTimeSeconds.Collect(ch)

	metrics.MetricRequestSizeBytes.Collect(ch)
	metrics.MetricRequestsTotal.Collect(ch)
	metrics.MetricRequestReadTimeSeconds.Collect(ch)
	metrics.MetricRequestReadUnmarshalTimeSeconds.Collect(ch)
	metrics.MetricRequestUnmarshalTimeSeconds.Collect(ch)

	metrics.MetricExternalResponseTimeSeconds.Collect(ch)

	metrics.MetricGRPCHandledTotal.Collect(ch)
	metrics.MetricGRPCReceivedTotal.Collect(ch)
	metrics.MetricGRPCSentTotal.Collect(ch)
	metrics.MetricGRPCStartedTotal.Collect(ch)

	delta := c.application.Uptime().Seconds() - c.upTime.Count()
	c.upTime.Add(delta)
	c.upTime.Collect(ch)

	c.buildInfo.With(
		"name", c.application.Name(),
		"version", c.application.Version(),
		"build", c.application.Build(),
		"build_date", c.application.BuildDate().Format(time.RFC3339),
		"start_date", c.application.StartDate().Format(time.RFC3339),
	).Set(1)
	c.buildInfo.Collect(ch)
}

func (c *Component) Metrics() snitch.Collector {
	return &metricsCollector{
		application: c.application,

		upTime:    snitch.NewCounter(metrics.MetricNameUpTime, "Application up time in seconds"),
		buildInfo: snitch.NewGauge(metrics.MetricNameBuildInfo, "Application build info"),
	}
}
