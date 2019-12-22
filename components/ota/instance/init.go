package instance

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/ota/internal"
)

func NewComponent() shadow.Component {
	return &internal.Component{}
}

func init() {
	shadow.MustRegisterComponent(NewComponent())
}
