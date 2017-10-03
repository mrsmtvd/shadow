package alerts

import (
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	Send(title string, message string, icon string)
	GetAlerts() []Alert
}
