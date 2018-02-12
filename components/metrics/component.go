package metrics

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/snitch"
)

type Component interface {
	shadow.Component

	Registry() snitch.Registerer
	Register(...snitch.Collector)
}
