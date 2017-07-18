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
	MetricWorkersStatus  = ComponentName + "_workers_status_total"
	MetricTasksTotal     = ComponentName + "_tasks_total"
	MetricTasksStatus    = ComponentName + "_tasks_status_total"
)

var (
	metricListenersTotal snitch.Gauge
	metricListenersTasks snitch.Gauge

	metricWorkersTotal        snitch.Gauge
	metricWorkerStatusWait    snitch.Gauge
	metricWorkerStatusProcess snitch.Gauge
	metricWorkerStatusBusy    snitch.Gauge

	metricTasksTotal               snitch.Gauge
	metricTasksStatusWait          snitch.Gauge
	metricTasksStatusProcess       snitch.Gauge
	metricTasksStatusSuccess       snitch.Gauge
	metricTasksStatusFail          snitch.Gauge
	metricTasksStatusFailByTimeout snitch.Gauge
	metricTasksStatusKill          snitch.Gauge
	metricTasksStatusRepeatWait    snitch.Gauge
)

type metricsCollector struct {
	dispatcher *dispatcher.Dispatcher
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	ch <- metricListenersTotal.Description()
	ch <- metricListenersTasks.Description()

	ch <- metricWorkersTotal.Description()
	ch <- metricWorkerStatusWait.Description()
	ch <- metricWorkerStatusProcess.Description()
	ch <- metricWorkerStatusBusy.Description()

	ch <- metricTasksTotal.Description()
	ch <- metricTasksStatusWait.Description()
	ch <- metricTasksStatusProcess.Description()
	ch <- metricTasksStatusSuccess.Description()
	ch <- metricTasksStatusFail.Description()
	ch <- metricTasksStatusFailByTimeout.Description()
	ch <- metricTasksStatusKill.Description()
	ch <- metricTasksStatusRepeatWait.Description()
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

	metricWorkerStatusWait.Set(workerStatusWait)
	metricWorkerStatusProcess.Set(workerStatusProcess)
	metricWorkerStatusBusy.Set(workerStatusBusy)

	ch <- metricListenersTotal
	ch <- metricListenersTasks

	ch <- metricWorkersTotal
	ch <- metricWorkerStatusWait
	ch <- metricWorkerStatusProcess
	ch <- metricWorkerStatusBusy

	ch <- metricTasksTotal
	ch <- metricTasksStatusWait
	ch <- metricTasksStatusProcess
	ch <- metricTasksStatusSuccess
	ch <- metricTasksStatusFail
	ch <- metricTasksStatusFailByTimeout
	ch <- metricTasksStatusKill
	ch <- metricTasksStatusRepeatWait
}

func (c *Component) Metrics() snitch.Collector {
	metricListenersTotal = snitch.NewGauge(MetricListenersTotal, "Number of listeners")
	metricListenersTasks = snitch.NewGauge(MetricListenersTasks, "Number of tasks in listeners")

	metricWorkersTotal = snitch.NewGauge(MetricWorkersTotal, "Number of workers")
	metricWorkerStatusWait = snitch.NewGauge(MetricWorkersStatus, "Number of workers with wait status", "status", "wait")
	metricWorkerStatusProcess = snitch.NewGauge(MetricWorkersStatus, "Number of workers with process status", "status", "process")
	metricWorkerStatusBusy = snitch.NewGauge(MetricWorkersStatus, "Number of workers with busy status", "status", "busy")

	metricTasksTotal = snitch.NewGauge(MetricTasksTotal, "Number of tasks")
	metricTasksStatusWait = snitch.NewGauge(MetricTasksStatus, "Number of tasks with wait status", "status", "wait")
	metricTasksStatusProcess = snitch.NewGauge(MetricTasksStatus, "Number of tasks with process status", "status", "process")
	metricTasksStatusSuccess = snitch.NewGauge(MetricTasksStatus, "Number of tasks with success status", "status", "success")
	metricTasksStatusFail = snitch.NewGauge(MetricTasksStatus, "Number of tasks with fail status", "status", "fail")
	metricTasksStatusFailByTimeout = snitch.NewGauge(MetricTasksStatus, "status", "fail-by-timeout")
	metricTasksStatusKill = snitch.NewGauge(MetricTasksStatus, "Number of tasks with kill status", "status", "kill")
	metricTasksStatusRepeatWait = snitch.NewGauge(MetricTasksStatus, "Number of tasks with wait repeat status", "status", "repeat-wait")

	return &metricsCollector{
		dispatcher: c.dispatcher,
	}
}
