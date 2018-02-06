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
	MetricListenersTotal       = workers.ComponentName + "_listeners_total"
	MetricListenersEventsTotal = workers.ComponentName + "_listeners_events_total"
	MetricWorkersTotal         = workers.ComponentName + "_workers_total"
	MetricWorkersLockedTotal   = workers.ComponentName + "_workers_locked_total"
	MetricTasksTotal           = workers.ComponentName + "_tasks_total"
	MetricTasksLockedTotal     = workers.ComponentName + "_tasks_locked_total"
)

var (
	metricListenersTotal       snitch.Gauge
	metricListenersEventsTotal snitch.Gauge
	metricWorkersTotal         snitch.Gauge
	metricWorkersLockedTotal   snitch.Gauge
	metricTasksTotal           snitch.Gauge
	metricTasksLockedTotal     snitch.Gauge
)

type metricsCollector struct {
	component *Component
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	metricListenersTotal.Describe(ch)
	metricListenersEventsTotal.Describe(ch)
	metricWorkersTotal.Describe(ch)
	metricWorkersLockedTotal.Describe(ch)
	metricTasksTotal.Describe(ch)
	metricTasksLockedTotal.Describe(ch)
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	workersLocked := 0
	workersList := c.component.GetWorkers()
	metricWorkersTotal.Set(float64(len(workersList)))

	for _, w := range workersList {
		if md := c.component.GetWorkerMetadata(w.Id()); md != nil {
			if md[ws.WorkerMetadataLocked].(bool) {
				workersLocked++
			}
		}
	}

	metricWorkersLockedTotal.Set(float64(workersLocked))

	tasksLocked := 0
	tasksList := c.component.GetTasks()
	metricTasksTotal.Set(float64(len(tasksList)))

	for _, t := range tasksList {
		if md := c.component.GetTaskMetadata(t.Id()); md != nil {
			if md[ws.TaskMetadataLocked].(bool) {
				tasksLocked++
			}
		}
	}

	metricTasksLockedTotal.Set(float64(tasksLocked))

	listeners := c.component.GetListeners()
	events := 0

	for _, l := range listeners {
		if md := c.component.GetListenerMetadata(l.Id()); md != nil {
			events += len(md[ws.ListenerMetadataEventIds].([]ws.EventId))
		}
	}

	metricListenersTotal.Set(float64(len(listeners)))
	metricListenersEventsTotal.Set(float64(events))

	metricListenersTotal.Collect(ch)
	metricListenersEventsTotal.Collect(ch)
	metricWorkersTotal.Collect(ch)
	metricWorkersLockedTotal.Collect(ch)
	metricTasksTotal.Collect(ch)
	metricTasksLockedTotal.Collect(ch)
}

func (c *metricsCollector) listener(_ context.Context, eventId ws.EventId, _ time.Time, args ...interface{}) {
	switch eventId {
	case ws.EventIdWorkerStatusChanged:
		metricWorkersTotal.With("status", strings.ToLower(args[2].(ws.Status).String())).Inc()
	case ws.EventIdTaskStatusChanged:
		metricTasksTotal.With("status", strings.ToLower(args[2].(ws.Status).String())).Inc()
	}
}

func (c *Component) Metrics() snitch.Collector {
	metricListenersTotal = snitch.NewGauge(MetricListenersTotal, "Number of listeners")
	metricListenersEventsTotal = snitch.NewGauge(MetricListenersEventsTotal, "Number of events of listeners")
	metricWorkersTotal = snitch.NewGauge(MetricWorkersTotal, "Number of workers")
	metricWorkersLockedTotal = snitch.NewGauge(MetricWorkersLockedTotal, "Number of locked workers")
	metricTasksTotal = snitch.NewGauge(MetricTasksTotal, "Number of tasks")
	metricTasksLockedTotal = snitch.NewGauge(MetricTasksLockedTotal, "Number of locked tasks")

	collector := &metricsCollector{
		component: c,
	}

	l := listener.NewFunctionListener(collector.listener)
	l.SetName(c.Name() + ".metrics")

	c.AddLockedListener(l.Id())
	c.AddListenerByEvents([]ws.EventId{
		ws.EventIdWorkerStatusChanged,
		ws.EventIdTaskStatusChanged,
	}, l)

	return collector
}
