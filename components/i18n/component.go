package i18n

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n/internationalization"
)

type Component interface {
	shadow.Component

	Manager() *internationalization.Manager
	LocaleFromRequest(*dashboard.Request) (*internationalization.Locale, error)
	LocaleFromAcceptLanguage(string) (*internationalization.Locale, error)
	LocaleFromSession(dashboard.Session) (*internationalization.Locale, error)
	SessionSave(dashboard.Session, *internationalization.Locale) error
}
