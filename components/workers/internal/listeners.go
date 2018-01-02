package internal

import (
	"context"
	"time"

	"github.com/kihamo/go-workers"
)

func (c *Component) listenerLogging(_ context.Context, eventId workers.EventId, _ time.Time, args ...interface{}) {
	switch eventId {
	case workers.EventIdWorkerAdd:
		id := args[0].(workers.Worker).Id()

		c.logger.Debugf("Added worker #%s", id, map[string]interface{}{
			"worker.id": id,
		})

	case workers.EventIdWorkerRemove:
		id := args[0].(workers.Worker).Id()

		c.logger.Debugf("Removed worker #%s", id, map[string]interface{}{
			"worker.id": id,
		})

	case workers.EventIdTaskAdd:
		id := args[0].(workers.Task).Id()

		c.logger.Debugf("Added task #%s", id, map[string]interface{}{
			"task.id": id,
		})

	case workers.EventIdTaskRemove:
		id := args[0].(workers.Task).Id()

		c.logger.Debugf("Removed task #%s", id, map[string]interface{}{
			"task.id": id,
		})

	case workers.EventIdListenerAdd:
		c.logger.Debug("Added listener", map[string]interface{}{
			"event":    args[0].(workers.EventId).String(),
			"listener": args[1].(workers.Listener),
		})

	case workers.EventIdListenerRemove:
		c.logger.Debug("Removed listener", map[string]interface{}{
			"event":    args[0].(workers.EventId).String(),
			"listener": args[1].(workers.Listener),
		})

	case workers.EventIdTaskExecuteStart:
		id := args[0].(workers.Task).Id()

		c.logger.Debugf("Execute task #%s started", id, map[string]interface{}{
			"task.id":   id,
			"worker.id": args[1].(workers.Worker).Id(),
		})

	case workers.EventIdTaskExecuteStop:
		fields := map[string]interface{}{
			"task.id":     args[0].(workers.Task).Id(),
			"worker.id":   args[1].(workers.Worker).Id(),
			"task.result": args[2],
		}

		if args[3] != nil {
			fields["task.err"] = args[3].(error).Error()
		}

		c.logger.Debugf("Execute task #%s stopped", fields["task.id"], fields)

	case workers.EventIdDispatcherStatusChanged:
		c.logger.Debug("Dispatcher status changed", map[string]interface{}{
			"dispatcher.status.current": args[1].(workers.Status).String(),
			"dispatcher.status.prev":    args[2].(workers.Status).String(),
		})

	case workers.EventIdWorkerStatusChanged:
		id := args[0].(workers.Worker).Id()

		c.logger.Debugf("Worker #%s status changed", id, map[string]interface{}{
			"worker.id":             id,
			"worker.status.current": args[1].(workers.Status).String(),
			"worker.status.prev":    args[2].(workers.Status).String(),
		})

	case workers.EventIdTaskStatusChanged:
		id := args[0].(workers.Task).Id()

		c.logger.Debugf("Task #%s status changed", id, map[string]interface{}{
			"task.id":             id,
			"task.status.current": args[1].(workers.Status).String(),
			"task.status.prev":    args[2].(workers.Status).String(),
		})
	}
}
