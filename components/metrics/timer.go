package metrics

import (
	"time"

	kit "github.com/go-kit/kit/metrics"
)

type Timer interface {
	With(...string) Timer
	ObserveDuration()
	ObserveDurationByTime(time.Time)
}

type MetricTimer struct {
	h kit.Histogram
	t time.Time
}

func NewMetricTimer(h kit.Histogram) Timer {
	return &MetricTimer{
		h: h,
		t: time.Now(),
	}
}

func (t *MetricTimer) With(labelValues ...string) Timer {
	t.h = t.h.With(labelValues...)
	return &MetricTimer{
		h: t.h.With(labelValues...),
		t: t.t,
	}
}

func (t *MetricTimer) ObserveDuration() {
	t.ObserveDurationByTime(t.t)
}

func (t *MetricTimer) ObserveDurationByTime(endTime time.Time) {
	d := time.Since(endTime).Seconds()
	if d < 0 {
		d = 0
	}
	t.h.Observe(d)
}
