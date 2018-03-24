package i18n

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n/internationalization"
)

func init() {
	dashboard.DefaultTemplateFunctions.AddFunction("i18n", templateFunctionTranslate)
	dashboard.DefaultTemplateFunctions.AddFunction("i18nPlural", templateFunctionTranslatePlural)
}

func templateFunctionTranslate(ID string, opts ...interface{}) string {
	return templateFunctionTranslatePlural(ID, "", 1, opts...)
}

func templateFunctionTranslatePlural(singleID, pluralID string, number int, _ ...interface{}) string {
	if number == 1 {
		return singleID
	}

	return pluralID
}

func NewOrNopFromRequest(request *dashboard.Request, application shadow.Application) *internationalization.Locale {
	if cmp := application.GetComponent(ComponentName); cmp != nil {
		locale, err := cmp.(Component).LocaleFromRequest(request)
		if err == nil {
			return locale
		}
	}

	return internationalization.NewLocale("nop")
}
