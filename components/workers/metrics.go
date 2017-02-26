package workers

import (
	kit "github.com/go-kit/kit/metrics"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/go-workers/worker"
	"github.com/kihamo/shadow/components/metrics"
)

const (
	MetricListenersTotal  = "workers.listeners.total"
	MetricListenersTasks  = "workers.listeners.tasks"
	MetricWorkersTotal  = "workers.workers.total"
	MetricWorkersStatus = "workers.workers.status"
	MetricTasksTotal    = "workers.tasks.total"
	MetricTasksStatus   = "workers.tasks.status"
)

var (
	metricListenersTotal kit.Gauge
	metricListenersTasks kit.Gauge
	metricWorkersTotal kit.Gauge
	metricTasksTotal   kit.Gauge

	metricWorkerStatusWait    kit.Gauge
	metricWorkerStatusProcess kit.Gauge
	metricWorkerStatusBusy    kit.Gauge

	metricTasksStatusWait          kit.Gauge
	metricTasksStatusProcess       kit.Gauge
	metricTasksStatusSuccess       kit.Gauge
	metricTasksStatusFail          kit.Gauge
	metricTasksStatusFailByTimeout kit.Gauge
	metricTasksStatusKill          kit.Gauge
	metricTasksStatusRepeatWait    kit.Gauge
)

func (c *Component) MetricsCapture() {
	metricListenersTotal.Set(float64(len(c.dispatcher.GetListeners())))
	metricListenersTasks.Set(float64(len(c.dispatcher.GetListenersTasks())))

	for _, w := range c.dispatcher.GetWorkers().GetItems() {
		metricWorkersTotal.Add(1)

		switch w.GetStatus() {
		case worker.WorkerStatusWait:
			metricWorkerStatusWait.Add(1)
		case worker.WorkerStatusProcess:
			metricWorkerStatusProcess.Add(1)
		case worker.WorkerStatusBusy:
			metricWorkerStatusBusy.Add(1)
		}
	}

	for _, t := range c.dispatcher.GetTasks().GetItems() {
		metricTasksTotal.Add(1)

		switch t.GetStatus() {
		case task.TaskStatusWait:
			metricTasksStatusWait.Add(1)
		case task.TaskStatusProcess:
			metricTasksStatusProcess.Add(1)
		case task.TaskStatusSuccess:
			metricTasksStatusSuccess.Add(1)
		case task.TaskStatusFail:
			metricTasksStatusFail.Add(1)
		case task.TaskStatusFailByTimeout:
			metricTasksStatusFailByTimeout.Add(1)
		case task.TaskStatusKill:
			metricTasksStatusKill.Add(1)
		case task.TaskStatusRepeatWait:
			metricTasksStatusRepeatWait.Add(1)
		}
	}
}

func (c *Component) MetricsRegister(m *metrics.Component) {
	metricListenersTotal = m.NewGauge(MetricListenersTotal)
	metricListenersTasks = m.NewGauge(MetricListenersTasks)

	metricWorkersTotal = m.NewGauge(MetricWorkersTotal)
	metricTasksTotal = m.NewGauge(MetricTasksTotal)

	metricWorkersStatus := m.NewGauge(MetricWorkersStatus)
	metricWorkerStatusWait = metricWorkersStatus.With("status", "wait")
	metricWorkerStatusProcess = metricWorkersStatus.With("status", "process")
	metricWorkerStatusBusy = metricWorkersStatus.With("status", "busy")

	metricTasksStatus := m.NewGauge(MetricTasksStatus)
	metricTasksStatusWait = metricTasksStatus.With("status", "wait")
	metricTasksStatusProcess = metricTasksStatus.With("status", "process")
	metricTasksStatusSuccess = metricTasksStatus.With("status", "success")
	metricTasksStatusFail = metricTasksStatus.With("status", "fail")
	metricTasksStatusFailByTimeout = metricTasksStatus.With("status", "fail-by-timeout")
	metricTasksStatusKill = metricTasksStatus.With("status", "kill")
	metricTasksStatusRepeatWait = metricTasksStatus.With("status", "repeat-wait")
}
