package workers

import (
	ws "github.com/kihamo/go-workers"
)

type ListenerWithEvents interface {
	ws.Listener

	Events() []ws.Event
}
