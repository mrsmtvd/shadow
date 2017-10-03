package instance

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/alerts/internal"
)

func NewComponent() shadow.Component {
	return &internal.Component{}
}
