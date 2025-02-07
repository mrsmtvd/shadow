package instance

import (
	"github.com/mrsmtvd/shadow"
	"github.com/mrsmtvd/shadow/components/metrics/internal"
)

func NewComponent() shadow.Component {
	return &internal.Component{}
}

func init() {
	shadow.MustRegisterComponent(NewComponent())
}
