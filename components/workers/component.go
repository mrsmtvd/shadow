package workers

import (
	ws "github.com/kihamo/go-workers"
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	AddSimpleWorker()
	AddWorker(ws.Worker)
	RemoveWorker(ws.Worker)
	GetWorkers() []ws.Worker
	GetWorkerMetadata(string) ws.Metadata

	AddTask(ws.Task)
	RemoveTask(ws.Task)
	GetTasks() []ws.Task
	GetTaskMetadata(string) ws.Metadata

	AddListener(ws.EventId, ws.Listener)
	RemoveListener(ws.EventId, ws.Listener)
	GetListeners() map[ws.EventId][]ws.Listener
}
