package metrics

import (
	"github.com/kihamo/snitch"
	"github.com/mrsmtvd/shadow"
)

type Component interface {
	shadow.Component

	Registry() snitch.Registerer
	Register(...snitch.Collector)
}
