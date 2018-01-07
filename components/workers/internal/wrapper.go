package internal

import (
	ws "github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/worker"
)

func (c *Component) AddSimpleWorker() {
	c.dispatcher.AddWorker(worker.NewSimpleWorker())
}

func (c *Component) AddWorker(worker ws.Worker) {
	c.dispatcher.AddWorker(worker)
}

func (c *Component) RemoveWorker(worker ws.Worker) {
	c.dispatcher.RemoveWorker(worker)
}

func (c *Component) GetWorkerMetadata(id string) ws.Metadata {
	return c.dispatcher.GetWorkerMetadata(id)
}

func (c *Component) GetWorkers() []ws.Worker {
	return c.dispatcher.GetWorkers()
}

func (c *Component) AddTask(task ws.Task) {
	c.dispatcher.AddTask(task)
}

func (c *Component) RemoveTask(task ws.Task) {
	c.dispatcher.RemoveTask(task)
}

func (c *Component) GetTaskMetadata(id string) ws.Metadata {
	return c.dispatcher.GetTaskMetadata(id)
}

func (c *Component) GetTasks() []ws.Task {
	return c.dispatcher.GetTasks()
}

func (c *Component) AddListenerByEvent(event ws.EventId, listener ws.Listener) {
	c.dispatcher.AddListener(event, listener)
}

func (c *Component) AddListenerByEvents(events []ws.EventId, listener ws.Listener) {
	for _, eventId := range events {
		c.AddListenerByEvent(eventId, listener)
	}
}

func (c *Component) RemoveListenerByEvent(event ws.EventId, listener ws.Listener) {
	c.dispatcher.RemoveListener(event, listener)
}

func (c *Component) RemoveListenerByEvents(events []ws.EventId, listener ws.Listener) {
	for _, eventId := range events {
		c.RemoveListenerByEvent(eventId, listener)
	}
}

func (c *Component) RemoveListener(listener ws.Listener) {
	md := c.GetListenerMetadata(listener.Id())

	if md == nil {
		return
	}

	mdValue, ok := md[ws.ListenerMetadataEventIds]
	if !ok {
		return
	}

	events, ok := mdValue.([]ws.EventId)
	if !ok {
		return
	}

	for _, eventId := range events {
		c.RemoveListenerByEvent(eventId, listener)
	}
}

func (c *Component) GetListenerMetadata(id string) ws.Metadata {
	return c.dispatcher.GetListenerMetadata(id)
}

func (c *Component) GetListeners() []ws.Listener {
	return c.dispatcher.GetListeners()
}
