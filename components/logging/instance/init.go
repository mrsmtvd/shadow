package instance

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/logging/internal"
)

func NewComponent() shadow.Component {
	return &internal.Component{}
}

func init() {
	shadow.MustRegisterComponent(NewComponent())
}
