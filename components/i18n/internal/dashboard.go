package internal

import (
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) DashboardTemplateFunctions() map[string]interface{} {
	return map[string]interface{}{
		"i18n":                c.templateFunctionTranslateContext,
		"i18nTranslate":       c.templateFunctionTranslate,
		"i18nPlural":          c.templateFunctionTranslatePluralContext,
		"i18nTranslatePlural": c.templateFunctionTranslatePlural,
	}
}

func (c *Component) templateFunctionTranslateContext(ID string, ctx map[string]interface{}, opts ...interface{}) string {
	return c.templateFunctionTranslatePluralContext(ID, "", 1, ctx, opts...)
}

func (c *Component) templateFunctionTranslate(ID string, opts ...interface{}) string {
	return c.templateFunctionTranslatePlural(ID, "", 1, opts...)
}

func (c *Component) templateFunctionTranslatePluralContext(singleID, pluralID string, number int, ctx map[string]interface{}, opts ...interface{}) string {
	if len(opts) <= 2 {
		if len(opts) <= 1 {
			if len(opts) == 0 {
				opts = append(opts, "")
			}

			componentName, ok := ctx["ComponentName"]
			if !ok {
				componentName = ""
			}
			opts = append(opts, componentName)
		}

		locale := ""
		if request_, ok_ := ctx["Request"]; ok_ {
			if request, ok := request_.(*dashboard.Request); ok {
				// TODO: in session

				// in request
				locale, err := c.LocaleFromAcceptLanguage(request.Original().Header.Get("Accept-Language"))
				if err == nil {
					opts = append(opts, locale.Locale())
				}
			}
		}
		opts = append(opts, locale)
	}

	return c.templateFunctionTranslatePlural(singleID, pluralID, number, opts...)
}

func (c *Component) templateFunctionTranslatePlural(singleID, pluralID string, number int, opts ...interface{}) string {
	var (
		context string // 0
		domain  string // 1
		locale  string // 2
	)

	if len(opts) > 0 {
		context = opts[0].(string)
	}

	if len(opts) > 1 {
		domain = opts[1].(string)
	}

	if len(opts) > 2 {
		locale = opts[2].(string)
	}

	return c.Manager().TranslatePlural(locale, domain, singleID, pluralID, number, context)
}
