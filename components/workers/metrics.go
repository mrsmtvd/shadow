package workers

import (
	kit "github.com/go-kit/kit/metrics"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/go-workers/worker"
	"github.com/kihamo/shadow/components/metrics"
)

const (
	MetricWorkersTotal  = "workers.workers.total"
	MetricWorkersStatus = "workers.workers.status"
	MetricTasksTotal    = "workers.tasks.total"
	MetricTasksStatus   = "workers.tasks.status"
)

var (
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
	workersTotal := 0
	tasksTotal := 0
	workersStatusWait := 0
	workersStatusProcess := 0
	workersStatusBusy := 0
	tasksStatusWait := 0
	tasksStatusProcess := 0
	tasksStatusSuccess := 0
	tasksStatusFail := 0
	tasksStatusFailByTimeout := 0
	tasksStatusKill := 0
	tasksStatusRepeatWait := 0

	for _, w := range c.dispatcher.GetWorkers().GetItems() {
		workersTotal += 1

		switch w.GetStatus() {
		case worker.WorkerStatusWait:
			workersStatusWait += 1
		case worker.WorkerStatusProcess:
			workersStatusProcess += 1
		case worker.WorkerStatusBusy:
			workersStatusBusy += 1
		}
	}

	for _, t := range c.dispatcher.GetTasks().GetItems() {
		tasksTotal += 1

		switch t.GetStatus() {
		case task.TaskStatusWait:
			tasksStatusWait += 1
		case task.TaskStatusProcess:
			tasksStatusProcess += 1
		case task.TaskStatusSuccess:
			tasksStatusSuccess += 1
		case task.TaskStatusFail:
			tasksStatusFail += 1
		case task.TaskStatusFailByTimeout:
			tasksStatusFailByTimeout += 1
		case task.TaskStatusKill:
			tasksStatusKill += 1
		case task.TaskStatusRepeatWait:
			tasksStatusRepeatWait += 1
		}
	}

	metricWorkersTotal.Set(float64(workersTotal))
	metricTasksTotal.Set(float64(tasksTotal))

	metricWorkerStatusWait.Set(float64(workersStatusWait))
	metricWorkerStatusProcess.Set(float64(workersStatusProcess))
	metricWorkerStatusBusy.Set(float64(workersStatusBusy))

	metricTasksStatusWait.Set(float64(tasksStatusWait))
	metricTasksStatusProcess.Set(float64(tasksStatusProcess))
	metricTasksStatusSuccess.Set(float64(tasksStatusSuccess))
	metricTasksStatusFail.Set(float64(tasksStatusFail))
	metricTasksStatusFailByTimeout.Set(float64(tasksStatusFailByTimeout))
	metricTasksStatusKill.Set(float64(tasksStatusKill))
	metricTasksStatusRepeatWait.Set(float64(tasksStatusRepeatWait))
}

func (c *Component) MetricsRegister(m *metrics.Component) {
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
