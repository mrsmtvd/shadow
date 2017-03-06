package workers

import (
	"sync"

	"github.com/kihamo/go-workers/dispatcher"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/go-workers/worker"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/metrics"
)

const (
	ComponentName = "workers"
)

type Component struct {
	application shadow.Application
	config      *config.Component
	logger      logger.Logger

	dispatcher *dispatcher.Dispatcher
}

func (c *Component) GetName() string {
	return ComponentName
}

func (c *Component) GetVersion() string {
	return "1.0.0"
}

func (c *Component) GetDependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: dashboard.ComponentName,
		},
		{
			Name: logger.ComponentName,
		},
		{
			Name: metrics.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	resourceConfig, err := a.GetComponent(config.ComponentName)
	if err != nil {
		return err
	}
	c.config = resourceConfig.(*config.Component)

	c.application = a

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) (err error) {
	c.logger = logger.NewOrNop(c.GetName(), c.application)

	c.dispatcher = dispatcher.NewDispatcher()
	c.setLogListener(wg)

	for i := 1; i <= c.config.GetInt(ConfigWorkersCount); i++ {
		c.AddWorker()
	}

	go func() {
		defer wg.Done()
		c.dispatcher.Run()

		c.dispatcher.SetTickerExecuteTasksDuration(c.config.GetDuration(ConfigTickerExecuteTasksDuration))
		c.dispatcher.SetTickerNotifyListenersDuration(c.config.GetDuration(ConfigTickerNotifyListenersDuration))
	}()

	return nil
}

func (c *Component) setLogListener(wg *sync.WaitGroup) {
	listener := dispatcher.NewDefaultListener()
	c.dispatcher.AddListener(listener)

	// logger for finished tasks
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case t := <-listener.GetTaskDoneChannel():
				if metricTasksTotal != nil {
					metricTasksTotal.Add(1)
				}

				status := t.GetStatus()

				switch status {
				case task.TaskStatusWait:
					c.logger.Debug("Finished", c.getLogFieldsForTask(t, map[string]interface{}{"task.status": "wait"}))

					if metricTasksStatusWait != nil {
						metricTasksStatusWait.Add(1)
					}
				case task.TaskStatusProcess:
					c.logger.Debug("Finished", c.getLogFieldsForTask(t, map[string]interface{}{"task.status": "process"}))

					if metricTasksStatusProcess != nil {
						metricTasksStatusProcess.Add(1)
					}
				case task.TaskStatusSuccess:
					c.logger.Debug("Success finished", c.getLogFieldsForTask(t, map[string]interface{}{"task.status": "success"}))

					if metricTasksStatusSuccess != nil {
						metricTasksStatusSuccess.Add(1)
					}
				case task.TaskStatusFail:
					c.logger.Error("Fail finished", c.getLogFieldsForTask(t, map[string]interface{}{"task.status": "fail"}))

					if metricTasksStatusFail != nil {
						metricTasksStatusFail.Add(1)
					}
				case task.TaskStatusFailByTimeout:
					c.logger.Error("Fail by timeout finished", c.getLogFieldsForTask(t, map[string]interface{}{"task.status": "fail-by-timeout"}))

					if metricTasksStatusFailByTimeout != nil {
						metricTasksStatusFailByTimeout.Add(1)
					}
				case task.TaskStatusKill:
					c.logger.Warn("Execute killed", c.getLogFieldsForTask(t, map[string]interface{}{"task.status": "kill"}))

					if metricTasksStatusKill != nil {
						metricTasksStatusKill.Add(1)
					}
				case task.TaskStatusRepeatWait:
					c.logger.Debug("Repeat execute", c.getLogFieldsForTask(t, map[string]interface{}{"task.status": "repeat-wait"}))

					if metricTasksStatusRepeatWait != nil {
						metricTasksStatusRepeatWait.Add(1)
					}
				default:
					c.logger.Warnf("Unknown task status %s", status, c.getLogFieldsForTask(t, map[string]interface{}{"task.status": "unknown"}))
				}
			}
		}
	}()
}

func (c *Component) AddTask(t task.Tasker) {
	c.logger.Debug("Add task", c.getLogFieldsForTask(t, nil))
	c.dispatcher.AddTask(t)
}

func (c *Component) AddNamedTaskByFunc(n string, f task.TaskFunction, a ...interface{}) task.Tasker {
	t := c.dispatcher.AddNamedTaskByFunc(n, f, a...)
	c.logger.Debug("Add task", c.getLogFieldsForTask(t, nil))
	return t
}

func (c *Component) AddTaskByFunc(f task.TaskFunction, a ...interface{}) task.Tasker {
	t := c.dispatcher.AddTaskByFunc(f, a...)
	c.logger.Debug("Add task", c.getLogFieldsForTask(t, nil))
	return t
}

func (c *Component) AddTaskByPriorityAndFunc(p int64, f task.TaskFunction, a ...interface{}) task.Tasker {
	t := c.dispatcher.AddTaskByPriorityAndFunc(p, f, a...)
	c.logger.Debug("Add task", c.getLogFieldsForTask(t, nil), map[string]interface{}{"priority": p})
	return t
}

func (c *Component) AddWorker() {
	w := c.dispatcher.AddWorker()
	c.logger.Debug("Add worker", map[string]interface{}{"worker.id": w.GetId()})
}

func (c *Component) GetWorkers() []worker.Worker {
	return c.dispatcher.GetWorkers().GetItems()
}

func (c *Component) getLogFieldsForTask(t task.Tasker, l map[string]interface{}) map[string]interface{} {
	fields := map[string]interface{}{
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

func (c *Component) AddListener(l dispatcher.Listener) {
	c.dispatcher.AddListener(l)
}

func (c *Component) RemoveListener(l dispatcher.Listener) {
	c.dispatcher.RemoveListener(l)
}
