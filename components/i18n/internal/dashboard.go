package internal

import (
	"context"
	"net/http"
	"reflect"
	"strconv"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/i18n"
	"github.com/mrsmtvd/shadow/components/i18n/internal/handlers"
	"github.com/mrsmtvd/shadow/components/logging"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	return []dashboard.Route{
		dashboard.RouteFromAssetFS(c),
		dashboard.NewRoute("/"+c.Name()+"/change/", handlers.NewChangeHandler(c)).
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

func (c *Component) DashboardMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// save locale in context
				if request := dashboard.RequestFromContext(r.Context()); request != nil {
					if locale, err := c.localeFromRequest(request); err == nil {
						request.WithContext(i18n.ContextWithLocale(r.Context(), locale))
						r = request.Original()
					}
				}

				next.ServeHTTP(w, r)
			})
		},
	}
}

func (c *Component) DashboardToolbar(ctx context.Context) string {
	locales := c.Manager().Locales()
	list := make([]string, 0, len(locales))

	for _, locale := range locales {
		list = append(list, locale.Locale())
	}

	content, err := c.dashboard.Renderer().RenderLayoutAndReturn(ctx, c.Name(), "toolbar", "blank", map[string]interface{}{
		"locales": list,
		"current": i18n.Locale(ctx).Locale(),
	})

	if err != nil {
		logging.Log(ctx).Error("Failed render toolbar", "error", err.Error())
	}

	return content
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

func (c *Component) templateFunctionTranslate(singleID string, opts ...interface{}) string {
	return c.templateFunctionTranslatePlural(singleID, "", 1, opts...)
}

func (c *Component) templateFunctionTranslatePlural(singleID, pluralID string, number interface{}, opts ...interface{}) string {
	var (
		ctx     map[string]interface{}
		callCtx string // 0
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
			callCtx = v
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
		if componentName, ok := ctx["NamespaceName"]; ok {
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
				locale = i18n.LocaleFromContext(request.Original().Context()).Locale()
			}
		}
	}

	return c.Manager().TranslatePlural(locale, domain, singleID, pluralID, c.convertToInt(number), callCtx, opts...)
}
