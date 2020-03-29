package internal

import (
	"context"
	"strings"
	"time"

	ws "github.com/kihamo/go-workers"
	"github.com/kihamo/shadow/components/metrics"
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
	metricListenersTotal       = snitch.NewGauge(MetricListenersTotal, "Number of listeners")
	metricListenersEventsTotal = snitch.NewGauge(MetricListenersEventsTotal, "Number of events of listeners")
	metricWorkersTotal         = snitch.NewGauge(MetricWorkersTotal, "Number of workers")
	metricWorkersLockedTotal   = snitch.NewGauge(MetricWorkersLockedTotal, "Number of locked workers")
	metricTasksTotal           = snitch.NewGauge(MetricTasksTotal, "Number of tasks")
	metricTasksLockedTotal     = snitch.NewGauge(MetricTasksLockedTotal, "Number of locked tasks")
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
			events += len(md[ws.ListenerMetadataEvents].([]ws.Event))
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

func (c *metricsCollector) listener(_ context.Context, event ws.Event, _ time.Time, args ...interface{}) {
	switch event {
	case ws.EventWorkerStatusChanged:
		metricWorkersTotal.With("status", strings.ToLower(args[2].(ws.Status).String())).Inc()
	case ws.EventTaskStatusChanged:
		metricTasksTotal.With("status", strings.ToLower(args[2].(ws.Status).String())).Inc()
	}
}

func (c *Component) Metrics() snitch.Collector {
	<-c.application.ReadyComponent(c.Name())

	collector := &metricsCollector{
		component: c,
	}

	l := NewListener(collector.listener, ws.EventWorkerStatusChanged, ws.EventTaskStatusChanged)
	l.SetName(c.Name() + "." + metrics.ComponentName)
	c.addLockedListener(l)

	return collector
}
