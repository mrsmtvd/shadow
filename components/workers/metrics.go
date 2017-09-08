package workers

import (
	"github.com/kihamo/go-workers/dispatcher"
	"github.com/kihamo/go-workers/worker"
	"github.com/kihamo/snitch"
)

const (
	MetricListenersTotal = ComponentName + "_listeners_total"
	MetricListenersTasks = ComponentName + "_listeners_tasks_total"
	MetricWorkersTotal   = ComponentName + "_workers_total"
	MetricTasksTotal     = ComponentName + "_tasks_total"
)

var (
	metricListenersTotal snitch.Gauge
	metricListenersTasks snitch.Gauge
	metricWorkersTotal   snitch.Gauge
	metricTasksTotal     snitch.Gauge
)

type metricsCollector struct {
	dispatcher *dispatcher.Dispatcher
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	metricListenersTotal.Describe(ch)
	metricListenersTasks.Describe(ch)
	metricWorkersTotal.Describe(ch)
	metricTasksTotal.Describe(ch)
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	metricListenersTotal.Set(float64(len(c.dispatcher.GetListeners())))
	metricListenersTasks.Set(float64(len(c.dispatcher.GetListenersTasks())))

	var (
		workerStatusWait    float64
		workerStatusProcess float64
		workerStatusBusy    float64
	)

	for _, w := range c.dispatcher.GetWorkers() {
		switch w.GetStatus() {
		case worker.WorkerStatusWait:
			workerStatusWait++
		case worker.WorkerStatusProcess:
			workerStatusProcess++
		case worker.WorkerStatusBusy:
			workerStatusBusy++
		}
	}

	metricWorkersTotal.With("status", "wait").Set(workerStatusWait)
	metricWorkersTotal.With("status", "process").Set(workerStatusProcess)
	metricWorkersTotal.With("status", "busy").Set(workerStatusBusy)

	metricListenersTotal.Collect(ch)
	metricListenersTasks.Collect(ch)
	metricWorkersTotal.Collect(ch)
	metricTasksTotal.Collect(ch)
}

func (c *Component) Metrics() snitch.Collector {
	metricListenersTotal = snitch.NewGauge(MetricListenersTotal, "Number of listeners")
	metricListenersTasks = snitch.NewGauge(MetricListenersTasks, "Number of tasks in listeners")
	metricWorkersTotal = snitch.NewGauge(MetricWorkersTotal, "Number of workers")
	metricTasksTotal = snitch.NewGauge(MetricTasksTotal, "Number of tasks")

	return &metricsCollector{
		dispatcher: c.dispatcher,
	}
}
