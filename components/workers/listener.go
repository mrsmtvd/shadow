package workers

import (
	ws "github.com/mrsmtvd/go-workers"
)

type ListenerWithEvents interface {
	ws.Listener

	Events() []ws.Event
}
