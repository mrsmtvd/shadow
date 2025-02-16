package internal

import (
	"context"
	"time"

	"github.com/mrsmtvd/go-workers"
	"github.com/mrsmtvd/go-workers/listener"
	"github.com/mrsmtvd/shadow/components/logging"
)

type Listener struct {
	listener.BaseListener

	function func(context.Context, workers.Event, time.Time, ...interface{})
	events   []workers.Event
}

func NewListener(function func(context.Context, workers.Event, time.Time, ...interface{}), events ...workers.Event) *Listener {
	l := &Listener{
		function: function,
		events:   events,
	}
	l.BaseListener.Init()

	return l
}

func (l *Listener) Run(ctx context.Context, event workers.Event, t time.Time, args ...interface{}) {
	l.function(ctx, event, t, args...)
}

func (l *Listener) Name() string {
	n := l.BaseListener.Name()

	if n == "" {
		return workers.FunctionName(l.function)
	}

	return n
}

func (l *Listener) Events() []workers.Event {
	return l.events
}

func (c *Component) newLoggingListener() *Listener {
	if !c.application.HasComponent(logging.ComponentName) {
		return nil
	}

	logNormal := c.logger.Debug
	logError := c.logger.Error

	l := NewListener(func(_ context.Context, event workers.Event, _ time.Time, args ...interface{}) {
		switch event {
		case workers.EventWorkerAdd:
			logNormal("Worker added", "worker.id", args[0].(workers.Worker).Id())

		case workers.EventWorkerRemove:
			logNormal("Worker removed", "worker.id", args[0].(workers.Worker).Id())

		case workers.EventTaskAdd:
			t := args[0].(workers.Task)

			logNormal("Task added", "task.id", t.Id(), "task.name", t.Name())

		case workers.EventTaskRemove:
			t := args[0].(workers.Task)

			logNormal("Task removed", "task.id", t.Id(), "task.name", t.Name())

		case workers.EventListenerAdd:
			e := args[0].(workers.Event)
			l := args[1].(workers.Listener)

			logNormal("Listener added", "event", e.String(), "listener.id", l.Id(), "listener.name", l.Name())

		case workers.EventListenerRemove:
			e := args[0].(workers.Event)
			l := args[1].(workers.Listener)

			logNormal("Listener removed", "event", e.String(), "listener.id", l.Id(), "listener.name", l.Name())

		case workers.EventTaskExecuteStart:
			t := args[0].(workers.Task)
			w := args[2].(workers.Worker)

			logNormal("Task execute started", "task.id", t.Id(), "task.name", t.Name(), "worker.id", w.Id())

		case workers.EventTaskExecuteStop:
			t := args[0].(workers.Task)
			w := args[2].(workers.Worker)

			fields := []interface{}{
				"task.id", t.Id(),
				"task.name", t.Name(),
				"worker.id", w.Id(),
				"task.result", args[4],
				"task.error", nil,
			}

			if args[5] != nil {
				fields = append(fields, "task.err", args[5].(error).Error())
				logError("Task ended with an error", fields...)
			}

			logNormal("Task execute stopped", fields...)

		case workers.EventDispatcherStatusChanged:
			logNormal("Dispatcher status changed",
				"dispatcher.status.current", args[1].(workers.Status).String(),
				"dispatcher.status.prev", args[2].(workers.Status).String(),
			)

		case workers.EventWorkerStatusChanged:
			w := args[0].(workers.Worker)

			logNormal("Worker status changed",
				"worker.id", w.Id(),
				"worker.status.current", args[2].(workers.Status).String(),
				"worker.status.prev", args[3].(workers.Status).String(),
			)

		case workers.EventTaskStatusChanged:
			t := args[0].(workers.Task)

			logNormal("Task status changed",
				"task.id", t.Id(),
				"task.name", t.Name(),
				"task.status.current", args[2].(workers.Status).String(),
				"task.status.prev", args[3].(workers.Status).String(),
			)

		default:
			logError("Fire unknown event",
				"event.id", event.Id(),
				"event.name", event.Name(),
				"args", args,
			)
		}
	}, workers.EventAll)

	l.SetName(c.Name() + "." + logging.ComponentName)

	return l
}
