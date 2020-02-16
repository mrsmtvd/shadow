package internal

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logging"
	m "github.com/kihamo/shadow/components/logging/http"
	"github.com/kihamo/shadow/components/logging/internal/handlers"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	return dashboard.NewMenu("Logging").
		WithUrl("/" + c.Name() + "/").
		WithIcon("headset")
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	levels := make([][]interface{}, 0)

	for _, v := range c.config.Variables() {
		if v.Key() == logging.ConfigLevel {
			opts := v.ViewOptions()
			if enum, ok := opts[config.ViewOptionEnumOptions]; ok {
				if s, ok := enum.([][]interface{}); ok {
					levels = s
				}
			}

			break
		}
	}

	return []dashboard.Route{
		dashboard.RouteFromAssetFS(c),
		dashboard.NewRoute("/"+c.Name()+"/", handlers.NewManagerHandler(c.global, levels)).
			WithMethods([]string{http.MethodGet, http.MethodPost}).
			WithAuth(true),
	}
}

func (c *Component) DashboardMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()

				// save logger in context
				if request := dashboard.RequestFromContext(ctx); request != nil {
					name := dashboard.TemplateNamespaceFromContext(ctx)
					request.WithContext(logging.ContextWithLogger(ctx, logging.NewLazyLogger(c.Logger(), name)))
					r = request.Original()
				}

				next.ServeHTTP(w, r)
			})
		},
		m.ServerMiddleware(c.Logger().Named(dashboard.ComponentName)),
	}
}
