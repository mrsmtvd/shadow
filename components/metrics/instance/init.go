package instance

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/metrics/internal"
)

func NewComponent() shadow.Component {
	return &internal.Component{}
}
