package internal

import (
	"strings"
	"sync"
	"time"

	"github.com/kihamo/snitch"
	jaeger "github.com/uber/jaeger-lib/metrics"
)

const (
	metricsScopeSeparator = "_"
)

var metricNormalizer = strings.NewReplacer(".", "_", "-", "_")

// TODO: необходимо кэшировать набор метрик, так как после переинициализации сздаются дубли метрик

func (c *Component) newMetricsFactory() jaeger.Factory {
	c.metricsOnce.Do(func() {
		c.metricsFactory = newFactory("", nil)
	})

	return c.metricsFactory
}

func (c *Component) Metrics() snitch.Collector {
	return c.newMetricsFactory().(snitch.Collector)
}

func newFactory(scope string, tags map[string]string) *factory {
	return &factory{
		scope:    scope,
		tags:     tags,
		children: make([]snitch.Collector, 0),
	}
}

type factory struct {
	snitch.Collector

	scope string
	tags  map[string]string

	mutex    sync.RWMutex
	children []snitch.Collector
}

func (f *factory) Describe(ch chan<- *snitch.Description) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	for _, m := range f.children {
		m.Describe(ch)
	}
}

func (f *factory) Collect(ch chan<- snitch.Metric) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	for _, m := range f.children {
		m.Collect(ch)
	}
}

func (f *factory) add(m snitch.Collector) {
	f.mutex.Lock()
	f.children = append(f.children, m)
	f.mutex.Unlock()
}

func (f *factory) Counter(opts jaeger.Options) jaeger.Counter {
	metric := snitch.NewCounter(f.subName(opts.Name), opts.Help, tagsToLabels(f.mergeTags(opts.Tags))...)
	f.add(metric)

	return &counter{
		m: metric,
	}
}

func (f *factory) Timer(opts jaeger.TimerOptions) jaeger.Timer {
	metric := snitch.NewTimer(f.subName(opts.Name), opts.Help, tagsToLabels(f.mergeTags(opts.Tags))...)
	f.add(metric)

	return &timer{
		m: metric,
	}
}

func (f *factory) Gauge(opts jaeger.Options) jaeger.Gauge {
	metric := snitch.NewGauge(f.subName(opts.Name), opts.Help, tagsToLabels(f.mergeTags(opts.Tags))...)
	f.add(metric)

	return &gauge{
		m: metric,
	}
}

func (f *factory) Histogram(opts jaeger.HistogramOptions) jaeger.Histogram {
	metric := snitch.NewHistogramWithQuantiles(f.subName(opts.Name), opts.Help, opts.Buckets, tagsToLabels(f.mergeTags(opts.Tags))...)
	f.add(metric)

	return &histogram{
		m: metric,
	}
}

func (f *factory) Namespace(scope jaeger.NSOptions) jaeger.Factory {
	metric := newFactory(scope.Name, scope.Tags)
	f.add(metric)

	return metric
}

func (f *factory) subName(name string) string {
	if f.scope == "" {
		return metricNormalizer.Replace(name)
	}

	if name == "" {
		return metricNormalizer.Replace(f.scope)
	}

	return metricNormalizer.Replace(f.scope + metricsScopeSeparator + name)
}

func (f *factory) mergeTags(tags map[string]string) map[string]string {
	result := make(map[string]string, len(f.tags)+len(tags))

	for k, v := range f.tags {
		result[k] = v
	}

	for k, v := range tags {
		result[k] = v
	}

	return result
}

type counter struct {
	m snitch.Counter
}

func (c *counter) Inc(value int64) {
	c.m.Add(float64(value))
}

type timer struct {
	m snitch.Timer
}

func (t *timer) Record(value time.Duration) {
	t.m.Update(value)
}

type gauge struct {
	m snitch.Gauge
}

func (g *gauge) Update(value int64) {
	g.m.Set(float64(value))
}

type histogram struct {
	m snitch.Histogram
}

func (h *histogram) Record(value float64) {
	h.m.Add(value)
}

func tagsToLabels(tags map[string]string) []string {
	labels := make([]string, 0, len(tags)*2)

	for k, v := range tags {
		labels = append(labels, k, v)
	}

	return labels
}
