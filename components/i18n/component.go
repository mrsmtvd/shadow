package i18n

import (
	"io"

	"github.com/mrsmtvd/shadow"
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/i18n/internationalization"
)

type Component interface {
	shadow.Component

	Manager() *internationalization.Manager
	LoadLocaleFromFiles(domain string, locales map[string][]io.ReadSeeker)
	SaveToSession(dashboard.Session, *internationalization.Locale)
	SaveToCookie(*dashboard.Response, *internationalization.Locale)
}
