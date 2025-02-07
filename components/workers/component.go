package workers

import (
	ws "github.com/kihamo/go-workers"
	"github.com/mrsmtvd/shadow"
)

type Component interface {
	shadow.Component

	AddSimpleWorker()
	LockedListeners() []ws.ListenerWithEvents

	AddWorker(ws.Worker)
	RemoveWorker(ws.Worker)
	GetWorkerMetadata(string) ws.Metadata
	GetWorkers() []ws.Worker

	AddTask(ws.Task)
	RemoveTask(ws.Task)
	GetTaskMetadata(string) ws.Metadata
	GetTasks() []ws.Task

	AddListener(ListenerWithEvents)
	AddListenerByEvent(ws.Event, ws.Listener)
	AddListenerByEvents([]ws.Event, ws.Listener)
	RemoveListenerByEvent(ws.Event, ws.Listener)
	RemoveListenerByEvents([]ws.Event, ws.Listener)
	RemoveListener(ws.Listener)
	GetListenerMetadata(string) ws.Metadata
	GetListeners() []ws.Listener
}
