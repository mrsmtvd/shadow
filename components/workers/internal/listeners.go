package internal

import (
	"time"

	"github.com/kihamo/go-workers"
)

func (c *Component) listenWorkerAdd(_ time.Time, args ...interface{}) {
	id := args[0].(workers.Worker).Id()

	c.logger.Debugf("Added worker #%s", id, map[string]interface{}{
		"worker.id": id,
	})
}

func (c *Component) listenWorkerRemove(_ time.Time, args ...interface{}) {
	id := args[0].(workers.Worker).Id()

	c.logger.Debugf("Removed worker #%s", id, map[string]interface{}{
		"worker.id": id,
	})
}

func (c *Component) listenTaskAdd(_ time.Time, args ...interface{}) {
	id := args[0].(workers.Task).Id()

	c.logger.Debugf("Added task #%s", id, map[string]interface{}{
		"task.id": id,
	})
}

func (c *Component) listenTaskRemove(_ time.Time, args ...interface{}) {
	id := args[0].(workers.Task).Id()

	c.logger.Debugf("Removed task #%s", id, map[string]interface{}{
		"task.id": id,
	})
}

func (c *Component) listenListenerAdd(_ time.Time, args ...interface{}) {
	c.logger.Debug("Added listener", map[string]interface{}{
		"event":    args[0].(workers.EventId).String(),
		"listener": args[1].(workers.Listener),
	})
}

func (c *Component) listenListenerRemove(_ time.Time, args ...interface{}) {
	c.logger.Debug("Removed listener", map[string]interface{}{
		"event":    args[0].(workers.EventId).String(),
		"listener": args[1].(workers.Listener),
	})
}

func (c *Component) listenTaskExecuteStart(_ time.Time, args ...interface{}) {
	id := args[0].(workers.Task).Id()

	c.logger.Debugf("Execute task #%s started", id, map[string]interface{}{
		"task.id":   id,
		"worker.id": args[1].(workers.Worker).Id(),
	})
}

func (c *Component) listenTaskExecuteStop(_ time.Time, args ...interface{}) {
	fields := map[string]interface{}{
		"task.id":     args[0].(workers.Task).Id(),
		"worker.id":   args[1].(workers.Worker).Id(),
		"task.result": args[2],
	}

	if args[3] != nil {
		fields["task.err"] = args[3].(error).Error()
	}

	c.logger.Debugf("Execute task #%s stopped", fields["task.id"], fields)
}

func (c *Component) listenDispatcherStatusChanged(_ time.Time, args ...interface{}) {
	c.logger.Debug("Dispatcher status changed", map[string]interface{}{
		"dispatcher.status.current": args[1].(workers.Status).String(),
		"dispatcher.status.prev":    args[2].(workers.Status).String(),
	})
}

func (c *Component) listenWorkerStatusChanged(_ time.Time, args ...interface{}) {
	id := args[0].(workers.Worker).Id()

	c.logger.Debugf("Worker #%s status changed", id, map[string]interface{}{
		"worker.id":             id,
		"worker.status.current": args[1].(workers.Status).String(),
		"worker.status.prev":    args[2].(workers.Status).String(),
	})
}

func (c *Component) listenTaskStatusChanged(_ time.Time, args ...interface{}) {
	id := args[0].(workers.Task).Id()

	c.logger.Debugf("Task #%s status changed", id, map[string]interface{}{
		"task.id":             id,
		"task.status.current": args[1].(workers.Status).String(),
		"task.status.prev":    args[2].(workers.Status).String(),
	})
}
