package metrics

import (
	"github.com/rcrowley/go-metrics"
)

func (r *Metrics) NewCounter(name string) metrics.Counter {
	return metrics.GetOrRegisterCounter(name, r.getRegistry())
}

func (r *Metrics) NewGauge(name string) metrics.Gauge {
	return metrics.GetOrRegisterGauge(name, r.getRegistry())
}

func (r *Metrics) NewGaugeFloat64(name string) metrics.GaugeFloat64 {
	return metrics.GetOrRegisterGaugeFloat64(name, r.getRegistry())
}

func (r *Metrics) NewGaugeHistogram(name string, sample metrics.Sample) metrics.Histogram {
	return metrics.GetOrRegisterHistogram(name, r.getRegistry(), sample)
}

func (r *Metrics) NewMeter(name string) metrics.Meter {
	return metrics.GetOrRegisterMeter(name, r.getRegistry())
}

func (r *Metrics) NewTimer(name string) metrics.Timer {
	return metrics.GetOrRegisterTimer(name, r.getRegistry())
}
