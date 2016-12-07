package metrics

import (
	"fmt"

	kit "github.com/go-kit/kit/metrics"
)

func (r *Resource) getName(name string) string {
	return fmt.Sprint(r.prefix, name)
}

func (r *Resource) NewCounter(name string) kit.Counter {
	return r.connector.NewCounter(r.getName(name))
}

func (r *Resource) NewGauge(name string) kit.Gauge {
	return r.connector.NewGauge(r.getName(name))
}

func (r *Resource) NewHistogram(name string) kit.Histogram {
	return r.connector.NewHistogram(r.getName(name))
}

func (r *Resource) NewTimer(name string) Timer {
	return NewMetricTimer(r.NewHistogram(name))
}
