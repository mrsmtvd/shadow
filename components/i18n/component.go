package i18n

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/i18n/internationalization"
)

type Component interface {
	shadow.Component

	Manager() *internationalization.Manager
}
