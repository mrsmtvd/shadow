package instance

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/examples/demo/components/demo/internal"
)

func NewComponent() shadow.Component {
	return &internal.Component{}
}

func init() {
	shadow.MustRegisterComponent(NewComponent())
}
