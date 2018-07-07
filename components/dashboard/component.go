package dashboard

import (
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	Renderer() Renderer
}
