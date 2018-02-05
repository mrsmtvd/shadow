package messengers

import (
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	RegisterMessenger(string, Messenger) error
	UnregisterMessenger(string)
	Messenger(string) Messenger
}
