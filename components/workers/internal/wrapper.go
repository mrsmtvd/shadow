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

func (c *Component) GetWorkers() []ws.Worker {
	return c.dispatcher.GetWorkers()
}

func (c *Component) GetWorkerMetadata(id string) ws.Metadata {
	return c.dispatcher.GetWorkerMetadata(id)
}

func (c *Component) AddTask(task ws.Task) {
	c.dispatcher.AddTask(task)
}

func (c *Component) RemoveTask(task ws.Task) {
	c.dispatcher.RemoveTask(task)
}

func (c *Component) GetTasks() []ws.Task {
	return c.dispatcher.GetTasks()
}

func (c *Component) GetTaskMetadata(id string) ws.Metadata {
	return c.dispatcher.GetTaskMetadata(id)
}

func (c *Component) AddListener(event ws.EventId, listener ws.Listener) {
	c.dispatcher.AddListener(event, listener)
}

func (c *Component) RemoveListener(event ws.EventId, listener ws.Listener) {
	c.dispatcher.RemoveListener(event, listener)
}

func (c *Component) GetListeners() map[ws.EventId][]ws.Listener {
	return c.dispatcher.GetListeners()
}
