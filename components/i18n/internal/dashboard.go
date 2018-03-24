package internal

import (
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) DashboardTemplateFunctions() map[string]interface{} {
	return map[string]interface{}{
		"i18n":       c.templateFunctionTranslate,
		"i18nPlural": c.templateFunctionTranslatePlural,
	}
}

func (c *Component) templateFunctionTranslate(ID string, opts ...interface{}) string {
	return c.templateFunctionTranslatePlural(ID, "", 1, opts...)
}

func (c *Component) templateFunctionTranslatePlural(singleID, pluralID string, number int, opts ...interface{}) string {
	var (
		ctx     map[string]interface{}
		context string // 0
		domain  string // 1
		locale  string // 2
	)

	// template context
	if len(opts) > 0 {
		if templateCtx, ok := opts[0].(map[string]interface{}); ok {
			ctx = templateCtx
			opts = opts[1:]
		}
	}

	// message context
	if len(opts) > 0 {
		if v, ok := opts[0].(string); ok {
			context = v
		}

		opts = opts[1:]
	}

	// message domain
	if len(opts) > 0 {
		if v, ok := opts[0].(string); ok {
			domain = v
		}

		opts = opts[1:]
	}

	if domain == "" && len(ctx) > 0 {
		if componentName, ok := ctx["ComponentName"]; ok {
			domain = componentName.(string)
		}
	}

	// message locale
	if len(opts) > 0 {
		if v, ok := opts[0].(string); ok {
			locale = v
		}

		opts = opts[1:]
	}

	if locale == "" && len(ctx) > 0 {
		if requestCtx, ok := ctx["Request"]; ok {
			if request, ok := requestCtx.(*dashboard.Request); ok {
				localeRequest, err := c.LocaleFromRequest(request)
				if err == nil {
					locale = localeRequest.Locale()
				}
			}
		}
	}

	return c.Manager().TranslatePlural(locale, domain, singleID, pluralID, number, context, opts...)
}
