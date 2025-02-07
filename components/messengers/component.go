package messengers

import (
	"github.com/mrsmtvd/shadow"
)

type Component interface {
	shadow.Component

	RegisterMessenger(string, Messenger) error
	UnregisterMessenger(string)
	Messenger(string) Messenger
}
