package i18n

import (
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n/internationalization"
)

func init() {
	dashboard.DefaultTemplateFunctions.AddFunction("i18n", templateFunctionTranslate)
	dashboard.DefaultTemplateFunctions.AddFunction("i18nPlural", templateFunctionTranslatePlural)
}

func templateFunctionTranslate(ID string, format ...interface{}) string {
	return templateFunctionTranslatePlural(ID, "", 1, format...)
}

func templateFunctionTranslatePlural(singleID, pluralID string, number int, format ...interface{}) string {
	if number == 1 {
		return internationalization.Format(singleID, format...)
	}

	return internationalization.Format(pluralID, format...)
}
