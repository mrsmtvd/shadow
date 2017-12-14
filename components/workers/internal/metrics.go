package internal

import (
	"strings"
	"time"

	w "github.com/kihamo/go-workers"
	"github.com/kihamo/shadow/components/workers"
	"github.com/kihamo/snitch"
)

const (
	MetricListenersTotal = workers.ComponentName + "_listeners_total"
	MetricWorkersTotal   = workers.ComponentName + "_workers_total"
	MetricTasksTotal     = workers.ComponentName + "_tasks_total"
)

var (
	metricListenersTotal snitch.Gauge
	metricWorkersTotal   snitch.Gauge
	metricTasksTotal     snitch.Gauge
)

type metricsCollector struct {
	component *Component
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	metricListenersTotal.Describe(ch)
	metricWorkersTotal.Describe(ch)
	metricTasksTotal.Describe(ch)
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {

	metricWorkersTotal.Set(float64(len(c.component.GetWorkers())))
	metricTasksTotal.Set(float64(len(c.component.GetTasks())))

	var totalListeners float64
	for _, list := range c.component.GetListeners() {
		totalListeners += float64(len(list))
	}
	metricListenersTotal.Set(totalListeners)

	metricListenersTotal.Collect(ch)
	metricWorkersTotal.Collect(ch)
	metricTasksTotal.Collect(ch)
}

func (c *metricsCollector) listenWorkerStatusChanged(_ time.Time, args ...interface{}) {
	metricWorkersTotal.With("status", strings.ToLower(args[1].(w.Status).String())).Inc()
}

func (c *metricsCollector) listenTaskStatusChanged(_ time.Time, args ...interface{}) {
	metricTasksTotal.With("status", strings.ToLower(args[1].(w.Status).String())).Inc()
}

func (c *Component) Metrics() snitch.Collector {
	metricListenersTotal = snitch.NewGauge(MetricListenersTotal, "Number of listeners")
	metricWorkersTotal = snitch.NewGauge(MetricWorkersTotal, "Number of workers")
	metricTasksTotal = snitch.NewGauge(MetricTasksTotal, "Number of tasks")

	collector := &metricsCollector{
		component: c,
	}

	collector.component.AddListener(w.EventIdWorkerStatusChanged, collector.listenWorkerStatusChanged)
	collector.component.AddListener(w.EventIdTaskStatusChanged, collector.listenTaskStatusChanged)

	return collector
}
