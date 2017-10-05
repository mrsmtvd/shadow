package workers

import (
	"github.com/kihamo/go-workers/dispatcher"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/go-workers/worker"
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	AddTask(t task.Tasker)
	AddNamedTaskByFunc(n string, f task.TaskFunction, a ...interface{}) task.Tasker
	AddTaskByFunc(f task.TaskFunction, a ...interface{}) task.Tasker
	AddTaskByPriorityAndFunc(p int64, f task.TaskFunction, a ...interface{}) task.Tasker
	RemoveTask(t task.Tasker)
	RemoveTaskById(id string)
	AddWorker()
	RemoveWorker(w worker.Worker)
	GetWorkers() []worker.Worker
	GetDefaultListenerName() string
	AddListener(l dispatcher.Listener)
	RemoveListener(l dispatcher.Listener)
}
