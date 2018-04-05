package internal

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n/internal/handlers"
)

func (c *Component) DashboardRoutes() []dashboard.Route {
	return []dashboard.Route{
		dashboard.NewRoute("/"+c.Name()+"/change/", &handlers.ChangeHandler{}).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
	}
}

func (c *Component) DashboardTemplateFunctions() map[string]interface{} {
	return map[string]interface{}{
		"i18n":       c.templateFunctionTranslate,
		"i18nPlural": c.templateFunctionTranslatePlural,
	}
}

func (c *Component) convertToInt(number interface{}) (cast int) {
	switch v := number.(type) {
	case string:
		cast, _ = strconv.Atoi(v)

	case bool:
		if v {
			cast = 1
		}

	case []byte:
		cast = c.convertToInt(string(v))

	case int64:
		cast = int(v)

	case int:
		cast = v

	case int8:
		cast = int(v)

	case int16:
		cast = int(v)

	case int32:
		cast = int(v)

	case uint:
		cast = int(v)

	case uint8:
		cast = int(v)

	case uint16:
		cast = int(v)

	case uint32:
		cast = int(v)

	case uint64:
		cast = int(v)

	case float32:
		cast = int(v)

	case float64:
		cast = int(v)

	default:
		if reflect.ValueOf(v).Kind() == reflect.Ptr {
			return c.convertToInt(reflect.Indirect(reflect.ValueOf(v)).Interface())
		}
	}

	return cast
}

func (c *Component) templateFunctionTranslate(ID string, opts ...interface{}) string {
	return c.templateFunctionTranslatePlural(ID, "", 1, opts...)
}

func (c *Component) templateFunctionTranslatePlural(singleID, pluralID string, number interface{}, opts ...interface{}) string {
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

	return c.Manager().TranslatePlural(locale, domain, singleID, pluralID, c.convertToInt(number), context, opts...)
}
