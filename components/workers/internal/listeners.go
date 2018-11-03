package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/kihamo/go-workers"
)

func (c *Component) listenerLogging(_ context.Context, event workers.Event, _ time.Time, args ...interface{}) {
	switch event {
	case workers.EventWorkerAdd:
		worker := args[0].(workers.Worker)

		c.logger.Debug(fmt.Sprintf("%s added", worker), "worker.id", worker.Id())

	case workers.EventWorkerRemove:
		worker := args[0].(workers.Worker)

		c.logger.Debug(fmt.Sprintf("%s removed", worker), "worker.id", worker.Id())

	case workers.EventTaskAdd:
		task := args[0].(workers.Task)

		c.logger.Debug(fmt.Sprintf("%s added", task),
			"task.id", task.Id(),
			"task.name", task.Name(),
		)

	case workers.EventTaskRemove:
		task := args[0].(workers.Task)

		c.logger.Debug(fmt.Sprintf("%s removed", task),
			"task.id", task.Id(),
			"task.name", task.Name(),
		)

	case workers.EventListenerAdd:
		event := args[0].(workers.Event)
		listener := args[1].(workers.Listener)

		c.logger.Debug(fmt.Sprintf("%s added for %s", listener, event),
			"event", event.String(),
			"listener.id", listener.Id(),
			"listener.name", listener.Name(),
		)

	case workers.EventListenerRemove:
		event := args[0].(workers.Event)
		listener := args[1].(workers.Listener)

		c.logger.Debug(fmt.Sprintf("%s removed for %s", listener, event),
			"event", event.String(),
			"listener.id", listener.Id(),
			"listener.name", listener.Name(),
		)

	case workers.EventTaskExecuteStart:
		task := args[0].(workers.Task)

		c.logger.Debug(fmt.Sprintf("%s execute started", task),
			"task.id", task.Id(),
			"task.name", task.Name(),
			"worker.id", args[2].(workers.Worker).Id(),
		)

	case workers.EventTaskExecuteStop:
		task := args[0].(workers.Task)

		fields := []interface{}{
			"task.id", task.Id(),
			"task.name", task.Name(),
			"worker.id", args[2].(workers.Worker).Id(),
			"task.result", args[4],
			"task.error", nil,
		}

		if args[5] != nil {
			fields = append(fields, "task.err", args[5].(error).Error())
			c.logger.Error("Task ended with an error", fields...)
		}

		c.logger.Debug(fmt.Sprintf("%s execute stopped", task), fields...)

	case workers.EventDispatcherStatusChanged:
		c.logger.Debug("Dispatcher status changed",
			"dispatcher.status.current", args[1].(workers.Status).String(),
			"dispatcher.status.prev", args[2].(workers.Status).String(),
		)

	case workers.EventWorkerStatusChanged:
		worker := args[0].(workers.Worker)

		c.logger.Debug(fmt.Sprintf("%s status changed", worker),
			"worker.id", worker.Id(),
			"worker.status.current", args[2].(workers.Status).String(),
			"worker.status.prev", args[3].(workers.Status).String(),
		)

	case workers.EventTaskStatusChanged:
		task := args[0].(workers.Task)

		c.logger.Debug(fmt.Sprintf("%s status changed", task),
			"task.id", task.Id(),
			"task.name", task.Name(),
			"task.status.current", args[2].(workers.Status).String(),
			"task.status.prev", args[3].(workers.Status).String(),
		)

	default:
		c.logger.Debug("Fire unknown event",
			"event.id", event.Id(),
			"event.name", event.Name(),
			"args", args,
		)
	}
}
