package workers

import (
	"sync"

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
	application *shadow.Application
	config      *config.Resource
	logger      xlog.Logger

	dispatcher *dispatcher.Dispatcher
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

	r.application = a

	return nil
}

func (r *Resource) Run(wg *sync.WaitGroup) (err error) {
	r.dispatcher = dispatcher.NewDispatcher()
	r.setLogListener(wg)

	if r.application.HasResource("metrics") {
		resourceMetrics, _ := r.application.GetResource("metrics")

		RegisterMetrics(resourceMetrics.(*metrics.Resource))
		metrics.CaptureMetrics(r.config.GetDuration("metrics.interval"), CaptureMetrics(r.dispatcher))
	}

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
					r.logger.Info("Finished", r.getLogFieldsForTask(t, map[string]interface{}{"task.status": "wait"}))
				case task.TaskStatusProcess:
					r.logger.Info("Finished", r.getLogFieldsForTask(t, map[string]interface{}{"task.status": "process"}))
				case task.TaskStatusSuccess:
					r.logger.Info("Success finished", r.getLogFieldsForTask(t, map[string]interface{}{"task.status": "success"}))
				case task.TaskStatusFail:
					r.logger.Error("Fail finished", r.getLogFieldsForTask(t, map[string]interface{}{"task.status": "fail"}))
				case task.TaskStatusFailByTimeout:
					r.logger.Error("Fail by timeout finished", r.getLogFieldsForTask(t, map[string]interface{}{"task.status": "fail-by-timeout"}))
				case task.TaskStatusKill:
					r.logger.Warn("Execute killed", r.getLogFieldsForTask(t, map[string]interface{}{"task.status": "kill"}))
				case task.TaskStatusRepeatWait:
					r.logger.Info("Repeat execute", r.getLogFieldsForTask(t, map[string]interface{}{"task.status": "repeat-wait"}))
				}
			}
		}
	}()
}

func (r *Resource) AddTask(t task.Tasker) {
	if r.logger != nil {
		r.logger.Info("Add task", r.getLogFieldsForTask(t, nil))
	}

	r.dispatcher.AddTask(t)
}

func (r *Resource) AddNamedTaskByFunc(n string, f task.TaskFunction, a ...interface{}) task.Tasker {
	t := r.dispatcher.AddNamedTaskByFunc(n, f, a...)

	if r.logger != nil {
		r.logger.Info("Add task", r.getLogFieldsForTask(t, nil))
	}

	return t
}

func (r *Resource) AddTaskByFunc(f task.TaskFunction, a ...interface{}) task.Tasker {
	t := r.dispatcher.AddTaskByFunc(f, a...)

	if r.logger != nil {
		r.logger.Info("Add task", r.getLogFieldsForTask(t, nil))
	}

	return t
}

func (r *Resource) AddWorker() {
	w := r.dispatcher.AddWorker()

	if r.logger != nil {
		r.logger.Infof("Add worker", xlog.F{"worker.id": w.GetId()})
	}
}

func (r *Resource) GetWorkers() []worker.Worker {
	return r.dispatcher.GetWorkers().GetItems()
}

func (r *Resource) getLogFieldsForTask(t task.Tasker, l map[string]interface{}) xlog.F {
	fields := xlog.F{
		"task.id":       t.GetId(),
		"task.function": t.GetFunctionName(),
		"task.priority": t.GetPriority(),
		"task.name":     t.GetName(),
		"task.duration": t.GetDuration().String(),
		"task.repeats":  t.GetRepeats(),
		"task.attemps":  t.GetAttempts(),
	}

	if lastError := t.GetLastError(); lastError != nil {
		fields["task.error"] = lastError
	}

	for k, v := range l {
		fields[k] = v
	}

	return fields
}

func (r *Resource) AddListener(l dispatcher.Listener) {
	r.dispatcher.AddListener(l)
}

func (r *Resource) RemoveListener(l dispatcher.Listener) {
	r.dispatcher.RemoveListener(l)
}
