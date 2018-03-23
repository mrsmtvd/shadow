package instance

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/i18n/internal"
)

func NewComponent() shadow.Component {
	return &internal.Component{}
}
