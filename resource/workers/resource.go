package workers

import (
	"sync"

	kitmetrics "github.com/go-kit/kit/metrics"
	"github.com/kihamo/go-workers/dispatcher"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/go-workers/worker"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
	"github.com/kihamo/shadow/resource/logger"
	"github.com/kihamo/shadow/resource/metrics"
	"github.com/rs/xlog"
)

type Resource struct {
	config  *config.Resource
	metrics *metrics.Resource
	logger  xlog.Logger

	dispatcher *dispatcher.Dispatcher

	metricWorkersTotal             kitmetrics.Counter
	metricTasksTotal               kitmetrics.Counter
	metricTasksStatusWait          kitmetrics.Counter
	metricTasksStatusProcess       kitmetrics.Counter
	metricTasksStatusSuccess       kitmetrics.Counter
	metricTasksStatusFail          kitmetrics.Counter
	metricTasksStatusFailByTimeout kitmetrics.Counter
	metricTasksStatusKill          kitmetrics.Counter
	metricTasksStatusRepeatWait    kitmetrics.Counter
}

func (r *Resource) GetName() string {
	return "workers"
}

func (r *Resource) Init(a *shadow.Application) error {
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	r.config = resourceConfig.(*config.Resource)

	if a.HasResource("logger") {
		resourceLogger, _ := a.GetResource("logger")
		r.logger = resourceLogger.(*logger.Resource).Get(r.GetName())
	}

	if a.HasResource("metrics") {
		resourceMetrics, _ := a.GetResource("metrics")
		r.metrics = resourceMetrics.(*metrics.Resource)
	}

	return nil
}

func (r *Resource) Run(wg *sync.WaitGroup) (err error) {
	if r.metrics != nil {
		r.metricWorkersTotal = r.metrics.NewCounter(MetricWorkersTotal)
		r.metricTasksTotal = r.metrics.NewCounter(MetricTasksTotal)

		metricTasksStatus := r.metrics.NewCounter(MetricTasksStatus)
		r.metricTasksStatusWait = metricTasksStatus.With("status", "wait")
		r.metricTasksStatusProcess = metricTasksStatus.With("status", "process")
		r.metricTasksStatusSuccess = metricTasksStatus.With("status", "success")
		r.metricTasksStatusFail = metricTasksStatus.With("status", "fail")
		r.metricTasksStatusFailByTimeout = metricTasksStatus.With("status", "fail-by-timeout")
		r.metricTasksStatusKill = metricTasksStatus.With("status", "kill")
		r.metricTasksStatusRepeatWait = metricTasksStatus.With("status", "repeat-wait")
	}

	r.dispatcher = dispatcher.NewDispatcher()
	r.setLogListener(wg)

	for i := 1; i <= r.config.GetInt("workers.count"); i++ {
		r.AddWorker()
	}

	go func() {
		defer wg.Done()
		r.dispatcher.Run()
	}()

	return nil
}

func (r *Resource) setLogListener(wg *sync.WaitGroup) {
	if r.logger == nil {
		return
	}

	listener := dispatcher.NewDefaultListener(r.config.GetInt("workers.count"))
	r.dispatcher.AddListener(listener)

	// logger for finished tasks
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case t := <-listener.TaskDone:
				switch t.GetStatus() {
				case task.TaskStatusWait:
					r.logger.Info("Finished", r.getLogFieldsForTask(t), xlog.F{"task.status": "wait"})

					if r.metricTasksStatusWait != nil {
						r.metricTasksStatusWait.Add(1)
					}
				case task.TaskStatusProcess:
					r.logger.Info("Finished", r.getLogFieldsForTask(t), xlog.F{"task.status": "process"})

					if r.metricTasksStatusProcess != nil {
						r.metricTasksStatusProcess.Add(1)
					}
				case task.TaskStatusSuccess:
					r.logger.Info("Success finished", r.getLogFieldsForTask(t), xlog.F{"task.status": "success"})

					if r.metricTasksStatusSuccess != nil {
						r.metricTasksStatusSuccess.Add(1)
					}
				case task.TaskStatusFail:
					r.logger.Error("Fail finished", r.getLogFieldsForTask(t), xlog.F{"task.status": "fail"})

					if r.metricTasksStatusFail != nil {
						r.metricTasksStatusFail.Add(1)
					}
				case task.TaskStatusFailByTimeout:
					r.logger.Error("Fail by timeout finished", r.getLogFieldsForTask(t), xlog.F{"task.status": "fail-by-timeout"})

					if r.metricTasksStatusFailByTimeout != nil {
						r.metricTasksStatusFailByTimeout.Add(1)
					}
				case task.TaskStatusKill:
					r.logger.Warn("Execute killed", r.getLogFieldsForTask(t), xlog.F{"task.status": "kill"})

					if r.metricTasksStatusKill != nil {
						r.metricTasksStatusKill.Add(1)
					}
				case task.TaskStatusRepeatWait:
					r.logger.Info("Repeat execute", r.getLogFieldsForTask(t), xlog.F{"task.status": "repeat-wait"})

					if r.metricTasksStatusRepeatWait != nil {
						r.metricTasksStatusRepeatWait.Add(1)
					}
				}
			}
		}
	}()
}

func (r *Resource) AddTask(t task.Tasker) {
	r.dispatcher.AddTask(t)

	if r.logger != nil {
		r.logger.Info("Add task", r.getLogFieldsForTask(t))
	}

	if r.metricTasksTotal != nil {
		r.metricTasksTotal.Add(1)
	}
}

func (r *Resource) AddNamedTaskByFunc(n string, f task.TaskFunction, a ...interface{}) task.Tasker {
	t := r.dispatcher.AddNamedTaskByFunc(n, f, a...)

	if r.logger != nil {
		r.logger.Info("Add task", r.getLogFieldsForTask(t))
	}

	if r.metricTasksTotal != nil {
		r.metricTasksTotal.Add(1)
	}

	return t
}

func (r *Resource) AddTaskByFunc(f task.TaskFunction, a ...interface{}) task.Tasker {
	t := r.dispatcher.AddTaskByFunc(f, a...)

	if r.logger != nil {
		r.logger.Info("Add task", r.getLogFieldsForTask(t))
	}

	if r.metricTasksTotal != nil {
		r.metricTasksTotal.Add(1)
	}

	return t
}

func (r *Resource) AddWorker() {
	w := r.dispatcher.AddWorker()

	if r.logger != nil {
		r.logger.Infof("Add worker", xlog.F{"worker.id": w.GetId()})
	}

	if r.metricWorkersTotal != nil {
		r.metricWorkersTotal.Add(1)
	}
}

func (r *Resource) GetWorkers() []worker.Worker {
	return r.dispatcher.GetWorkers().GetItems()
}

func (r *Resource) getLogFieldsForTask(t task.Tasker) xlog.F {
	fields := xlog.F{
		"task.id":        t.GetId(),
		"task.function":  t.GetFunctionName(),
		"task.arguments": t.GetArguments(),
		"task.priority":  t.GetPriority(),
		"task.name":      t.GetName(),
		"task.duration":  t.GetDuration(),
		"task.repeats":   t.GetRepeats(),
		"task.attemps":   t.GetAttempts(),
	}

	if lastError := t.GetLastError(); lastError != nil {
		fields["task.error"] = lastError
	}

	return fields
}

func (r *Resource) AddListener(l dispatcher.Listener) {
	r.dispatcher.AddListener(l)
}

func (r *Resource) RemoveListener(l dispatcher.Listener) {
	r.dispatcher.RemoveListener(l)
}
