package internal

import (
	ws "github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/worker"
	"github.com/mrsmtvd/shadow/components/workers"
)

func (c *Component) AddSimpleWorker() {
	_ = c.dispatcher.AddWorker(worker.NewSimpleWorker())
}

func (c *Component) AddWorker(worker ws.Worker) {
	_ = c.dispatcher.AddWorker(worker)
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
	_ = c.dispatcher.AddTask(task)
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

func (c *Component) AddListener(listener workers.ListenerWithEvents) {
	c.AddListenerByEvents(listener.Events(), listener)
}

func (c *Component) AddListenerByEvent(event ws.Event, listener ws.Listener) {
	_ = c.dispatcher.AddListener(event, listener)
}

func (c *Component) AddListenerByEvents(events []ws.Event, listener ws.Listener) {
	for _, event := range events {
		c.AddListenerByEvent(event, listener)
	}
}

func (c *Component) RemoveListenerByEvent(event ws.Event, listener ws.Listener) {
	c.dispatcher.RemoveListener(event, listener)
}

func (c *Component) RemoveListenerByEvents(events []ws.Event, listener ws.Listener) {
	for _, event := range events {
		c.RemoveListenerByEvent(event, listener)
	}
}

func (c *Component) RemoveListener(listener ws.Listener) {
	md := c.GetListenerMetadata(listener.Id())

	if md == nil {
		return
	}

	mdValue, ok := md[ws.ListenerMetadataEvents]
	if !ok {
		return
	}

	events, ok := mdValue.([]ws.Event)
	if !ok {
		return
	}

	for _, event := range events {
		c.RemoveListenerByEvent(event, listener)
	}
}

func (c *Component) GetListenerMetadata(id string) ws.Metadata {
	return c.dispatcher.GetListenerMetadata(id)
}

func (c *Component) GetListeners() []ws.Listener {
	return c.dispatcher.GetListeners()
}
