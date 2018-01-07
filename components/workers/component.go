package workers

import (
	ws "github.com/kihamo/go-workers"
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	AddSimpleWorker()
	GetLockedListeners() []string
	AddLockedListener(string)

	AddWorker(ws.Worker)
	RemoveWorker(ws.Worker)
	GetWorkerMetadata(string) ws.Metadata
	GetWorkers() []ws.Worker

	AddTask(ws.Task)
	RemoveTask(ws.Task)
	GetTaskMetadata(string) ws.Metadata
	GetTasks() []ws.Task

	AddListenerByEvent(ws.EventId, ws.Listener)
	AddListenerByEvents([]ws.EventId, ws.Listener)
	RemoveListenerByEvent(ws.EventId, ws.Listener)
	RemoveListenerByEvents([]ws.EventId, ws.Listener)
	RemoveListener(ws.Listener)
	GetListenerMetadata(string) ws.Metadata
	GetListeners() []ws.Listener
}
