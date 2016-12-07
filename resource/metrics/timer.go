package metrics

import (
	"time"

	kit "github.com/go-kit/kit/metrics"
)

type Timer struct {
	h kit.Histogram
	t time.Time
}

func NewTimer(h kit.Histogram) *Timer {
	return &Timer{
		h: h,
		t: time.Now(),
	}
}

func (t *Timer) With(labelValues ...string) *Timer {
	t.h = t.h.With(labelValues...)
	return &Timer{
		h: t.h.With(labelValues...),
		t: t.t,
	}
}

func (t *Timer) ObserveDuration() {
	d := time.Since(t.t).Seconds()
	if d < 0 {
		d = 0
	}
	t.h.Observe(d)
}
