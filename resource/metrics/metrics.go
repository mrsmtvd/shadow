package metrics

import (
	"github.com/go-kit/kit/metrics/influx"
)

func (r *Metrics) NewCounter(name string) *influx.Counter {
	return r.connector.NewCounter(name)
}

func (r *Metrics) NewGauge(name string) *influx.Gauge {
	return r.connector.NewGauge(name)
}

func (r *Metrics) NewHistogram(name string) *influx.Histogram {
	return r.connector.NewHistogram(name)
}

func (r *Metrics) NewTimer(name string) *Timer {
	return NewTimer(r.NewHistogram(name))
}
