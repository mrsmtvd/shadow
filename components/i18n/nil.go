package i18n

import (
	"github.com/kihamo/shadow/components/dashboard"
)

func init() {
	dashboard.DefaultTemplateFunctions.AddFunction("i18n", templateFunctionTranslateContext)
	dashboard.DefaultTemplateFunctions.AddFunction("i18nTranslate", templateFunctionTranslate)
	dashboard.DefaultTemplateFunctions.AddFunction("i18nPlural", templateFunctionTranslatePluralContext)
	dashboard.DefaultTemplateFunctions.AddFunction("i18nTranslatePlural", templateFunctionTranslatePlural)
}

func templateFunctionTranslateContext(ID string, ctx map[string]interface{}, opts ...interface{}) string {
	return templateFunctionTranslatePluralContext(ID, "", 1, ctx, opts...)
}

func templateFunctionTranslate(ID string, opts ...interface{}) string {
	return templateFunctionTranslatePlural(ID, "", 1, opts...)
}

func templateFunctionTranslatePluralContext(singleID, pluralID string, number int, ctx map[string]interface{}, opts ...interface{}) string {
	return templateFunctionTranslatePlural(singleID, pluralID, number, opts...)
}

func templateFunctionTranslatePlural(singleID, pluralID string, number int, _ ...interface{}) string {
	if number == 1 {
		return singleID
	}

	return pluralID
}
