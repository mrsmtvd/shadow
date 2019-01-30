package i18n

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n/internationalization"
)

type Component interface {
	shadow.Component

	Manager() *internationalization.Manager
	SaveToSession(dashboard.Session, *internationalization.Locale) error
	SaveToCookie(*dashboard.Response, *internationalization.Locale) error
}
