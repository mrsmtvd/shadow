package workers

const (
	MetricWorkersInWaitStatus          = "workers.status.wait"
	MetricWorkersInProccessStatus      = "workers.status.proccess"
	MetricWorkersInSuccessStatus       = "workers.status.success"
	MetricWorkersInFailStatus          = "workers.status.fail"
	MetricWorkersInFailByTimeOutStatus = "workers.status.fail_by_timeout"
	MetricWorkersInKillStatus          = "workers.status.kill"
	MetricWorkersInRepeatWaitStatus    = "workers.status.repeat_wait"

	MetricTotalWorkers = "workers.total.workers"
	MetricTotalTasks   = "workers.total.tasks"
)
