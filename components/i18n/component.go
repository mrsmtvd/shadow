package i18n

import (
	"io"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n/internationalization"
)

type Component interface {
	shadow.Component

	Manager() *internationalization.Manager
	LoadLocaleFromFiles(domain string, locales map[string][]io.ReadSeeker)
	SaveToSession(dashboard.Session, *internationalization.Locale)
	SaveToCookie(*dashboard.Response, *internationalization.Locale)
}
