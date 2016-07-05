package resource

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/kihamo/go-workers/dispatcher"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/go-workers/worker"
	"github.com/kihamo/shadow"
)

type Workers struct {
	config     *Config
	logger     *logrus.Entry
	finish     chan task.Tasker
	dispatcher *dispatcher.Dispatcher
}

func (r *Workers) GetName() string {
	return "workers"
}

func (r *Workers) GetConfigVariables() []ConfigVariable {
	return []ConfigVariable{
		ConfigVariable{
			Key:   "workers.count",
			Value: 2,
			Usage: "Default workers count",
		},
		ConfigVariable{
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
	r.config = resourceConfig.(*Config)

	resourceLogger, err := a.GetResource("logger")
	if err != nil {
		return err
	}
	r.logger = resourceLogger.(*Logger).Get(r.GetName())

	return nil
}

func (r *Workers) Run(wg *sync.WaitGroup) (err error) {
	r.finish = make(chan task.Tasker, r.config.GetInt64("workers.done.size"))

	r.dispatcher = dispatcher.NewDispatcher()
	r.dispatcher.SetTaskDoneChannel(r.finish)

	for i := 1; i <= int(r.config.GetInt64("workers.count")); i++ {
		r.AddWorker()
	}

	// logger for finished tasks
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case t := <-r.finish:
				switch t.GetStatus() {
				case task.TaskStatusWait:
					r.getLogEntryForTask(t).
						WithField("task.status", "wait").
						Info("Finished")
				case task.TaskStatusProcess:
					r.getLogEntryForTask(t).
						WithField("task.status", "process").
						Info("Finished")
				case task.TaskStatusSuccess:
					r.getLogEntryForTask(t).
						WithField("task.status", "success").
						Info("Success finished")
				case task.TaskStatusFail:
					r.getLogEntryForTask(t).
						WithField("task.status", "fail").
						Error("Fail finished")
				case task.TaskStatusFailByTimeout:
					r.getLogEntryForTask(t).
						WithField("task.status", "fail-by-timeout").
						Error("Fail by timeout finished")
				case task.TaskStatusKill:
					r.getLogEntryForTask(t).
						WithField("task.status", "fail-by-timeout").
						Warn("Execute killed")
				case task.TaskStatusRepeatWait:
					r.getLogEntryForTask(t).
						WithField("task.status", "repeat-wait").
						Info("Repeat execute")
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		r.dispatcher.Run()
	}()

	return nil
}

func (r *Workers) AddTask(t task.Tasker) {
	r.dispatcher.AddTask(t)

	r.getLogEntryForTask(t).Info("Add task")
}

func (r *Workers) AddNamedTaskByFunc(n string, f task.TaskFunction, a ...interface{}) task.Tasker {
	t := r.dispatcher.AddNamedTaskByFunc(n, f, a...)

	r.getLogEntryForTask(t).Info("Add task")

	return t
}

func (r *Workers) AddTaskByFunc(f task.TaskFunction, a ...interface{}) task.Tasker {
	t := r.dispatcher.AddTaskByFunc(f, a...)

	r.getLogEntryForTask(t).Info("Add task")

	return t
}

func (r *Workers) AddWorker() {
	w := r.dispatcher.AddWorker()

	r.logger.WithField("worker.id", w.GetId()).Info("Add worker")
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
