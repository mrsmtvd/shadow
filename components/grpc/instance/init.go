package instance

import (
	"github.com/mrsmtvd/shadow"
	"github.com/mrsmtvd/shadow/components/grpc/internal"
)

func NewComponent() shadow.Component {
	return &internal.Component{}
}

func init() {
	shadow.MustRegisterComponent(NewComponent())
}
