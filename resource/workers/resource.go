package workers

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/kihamo/go-workers/dispatcher"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/go-workers/worker"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
	"github.com/kihamo/shadow/resource/logger"
	"github.com/kihamo/shadow/resource/metrics"
)

type Workers struct {
	config     *config.Config
	logger     *logrus.Entry
	metrics    *metrics.Metrics
	dispatcher *dispatcher.Dispatcher
}

func (r *Workers) GetName() string {
	return "workers"
}

func (r *Workers) GetConfigVariables() []config.ConfigVariable {
	return []config.ConfigVariable{
		config.ConfigVariable{
			Key:   "workers.count",
			Value: 2,
			Usage: "Default workers count",
		},
		config.ConfigVariable{
			Key:   "workers.done.size",
			Value: 1000,
			Usage: "Size buffer of done task channel",
		},
	}
}

func (r *Workers) Init(a *shadow.Application) error {
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	r.config = resourceConfig.(*config.Config)

	if a.HasResource("logger") {
		resourceLogger, _ := a.GetResource("logger")
		r.logger = resourceLogger.(*logger.Logger).Get(r.GetName())
	}

	if a.HasResource("metrics") {
		resourceMetrics, _ := a.GetResource("metrics")
		r.metrics = resourceMetrics.(*metrics.Metrics)
	}

	return nil
}

func (r *Workers) Run(wg *sync.WaitGroup) (err error) {
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

func (r *Workers) setLogListener(wg *sync.WaitGroup) {
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
					r.getLogEntryForTask(t).
						WithField("task.status", "wait").
						Info("Finished")

					if r.metrics != nil {
						r.metrics.NewCounter(MetricWorkersInWaitStatus).Inc(1)
					}
				case task.TaskStatusProcess:
					r.getLogEntryForTask(t).
						WithField("task.status", "process").
						Info("Finished")

					if r.metrics != nil {
						r.metrics.NewCounter(MetricWorkersInProccessStatus).Inc(1)
					}
				case task.TaskStatusSuccess:
					r.getLogEntryForTask(t).
						WithField("task.status", "success").
						Info("Success finished")

					if r.metrics != nil {
						r.metrics.NewCounter(MetricWorkersInSuccessStatus).Inc(1)
					}
				case task.TaskStatusFail:
					r.getLogEntryForTask(t).
						WithField("task.status", "fail").
						Error("Fail finished")

					if r.metrics != nil {
						r.metrics.NewCounter(MetricWorkersInFailStatus).Inc(1)
					}
				case task.TaskStatusFailByTimeout:
					r.getLogEntryForTask(t).
						WithField("task.status", "fail-by-timeout").
						Error("Fail by timeout finished")

					if r.metrics != nil {
						r.metrics.NewCounter(MetricWorkersInFailByTimeOutStatus).Inc(1)
					}
				case task.TaskStatusKill:
					r.getLogEntryForTask(t).
						WithField("task.status", "kill").
						Warn("Execute killed")

					if r.metrics != nil {
						r.metrics.NewCounter(MetricWorkersInKillStatus).Inc(1)
					}
				case task.TaskStatusRepeatWait:
					r.getLogEntryForTask(t).
						WithField("task.status", "repeat-wait").
						Info("Repeat execute")

					if r.metrics != nil {
						r.metrics.NewCounter(MetricWorkersInRepeatWaitStatus).Inc(1)
					}
				}
			}
		}
	}()
}

func (r *Workers) AddTask(t task.Tasker) {
	r.dispatcher.AddTask(t)

	if r.logger != nil {
		r.getLogEntryForTask(t).Info("Add task")
	}

	if r.metrics != nil {
		r.metrics.NewCounter(MetricTotalTasks).Inc(1)
	}
}

func (r *Workers) AddNamedTaskByFunc(n string, f task.TaskFunction, a ...interface{}) task.Tasker {
	t := r.dispatcher.AddNamedTaskByFunc(n, f, a...)

	if r.logger != nil {
		r.getLogEntryForTask(t).Info("Add task")
	}

	if r.metrics != nil {
		r.metrics.NewCounter(MetricTotalTasks).Inc(1)
	}

	return t
}

func (r *Workers) AddTaskByFunc(f task.TaskFunction, a ...interface{}) task.Tasker {
	t := r.dispatcher.AddTaskByFunc(f, a...)

	if r.logger != nil {
		r.getLogEntryForTask(t).Info("Add task")
	}

	if r.metrics != nil {
		r.metrics.NewCounter(MetricTotalTasks).Inc(1)
	}

	return t
}

func (r *Workers) AddWorker() {
	w := r.dispatcher.AddWorker()

	if r.logger != nil {
		r.logger.WithField("worker.id", w.GetId()).Info("Add worker")
	}

	if r.metrics != nil {
		r.metrics.NewCounter(MetricTotalWorkers).Inc(1)
	}
}

func (r *Workers) GetWorkers() []worker.Worker {
	return r.dispatcher.GetWorkers().GetItems()
}

func (r *Workers) getLogEntryForTask(t task.Tasker) *logrus.Entry {
	entry := r.logger.WithFields(logrus.Fields{
		"task.id":        t.GetId(),
		"task.function":  t.GetFunctionName(),
		"task.arguments": t.GetArguments(),
		"task.priority":  t.GetPriority(),
		"task.name":      t.GetName(),
		"task.duration":  t.GetDuration(),
		"task.repeats":   t.GetRepeats(),
		"task.attemps":   t.GetAttempts(),
	})

	lastError := t.GetLastError()
	if lastError != nil {
		entry = entry.WithField("task.error", lastError)
	}

	return entry
}

func (r *Workers) AddListener(l dispatcher.Listener) {
	r.dispatcher.AddListener(l)
}

func (r *Workers) RemoveListener(l dispatcher.Listener) {
	r.dispatcher.RemoveListener(l)
}
