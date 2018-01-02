package internal

import (
	"context"
	"strings"
	"time"

	ws "github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
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
	metricListenersTotal.Set(float64(len(c.component.GetListeners())))

	metricListenersTotal.Collect(ch)
	metricWorkersTotal.Collect(ch)
	metricTasksTotal.Collect(ch)
}

func (c *metricsCollector) listener(_ context.Context, eventId ws.EventId, _ time.Time, args ...interface{}) {
	switch eventId {
	case ws.EventIdWorkerStatusChanged:
		metricWorkersTotal.With("status", strings.ToLower(args[1].(ws.Status).String())).Inc()
	case ws.EventIdTaskStatusChanged:
		metricTasksTotal.With("status", strings.ToLower(args[1].(ws.Status).String())).Inc()
	}
}

func (c *Component) Metrics() snitch.Collector {
	metricListenersTotal = snitch.NewGauge(MetricListenersTotal, "Number of listeners")
	metricWorkersTotal = snitch.NewGauge(MetricWorkersTotal, "Number of workers")
	metricTasksTotal = snitch.NewGauge(MetricTasksTotal, "Number of tasks")

	collector := &metricsCollector{
		component: c,
	}

	l := listener.NewFunctionListener(collector.listener)
	l.SetName(c.GetName() + ".metrics")

	c.AddLockedListener(l.Id())
	c.AddListener(ws.EventIdWorkerStatusChanged, l)
	c.AddListener(ws.EventIdTaskStatusChanged, l)

	return collector
}
