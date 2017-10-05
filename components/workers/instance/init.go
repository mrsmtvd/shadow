package instance

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/workers/internal"
)

func NewComponent() shadow.Component {
	return &internal.Component{}
}
