package internal

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/config/internal/handlers"
	"github.com/mrsmtvd/shadow/components/dashboard"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	return dashboard.NewMenu("Configuration").
		WithURL("/" + c.Name() + "/").
		WithIcon("cog")
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	return []dashboard.Route{
		dashboard.RouteFromAssetFS(c),
		dashboard.NewRoute("/"+c.Name()+"/", handlers.NewManagerHandler(c)).
			WithMethods([]string{http.MethodGet, http.MethodPost}).
			WithAuth(true),
	}
}

func (c *Component) DashboardTemplateFunctions() map[string]interface{} {
	return map[string]interface{}{
		"config": c.templateFunctionConfig,
	}
}

func (c *Component) DashboardMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// save config in context
				if request := dashboard.RequestFromContext(r.Context()); request != nil {
					request.WithContext(config.ContextWithConfig(r.Context(), c))
					r = request.Original()
				}

				next.ServeHTTP(w, r)
			})
		},
	}
}

func (c *Component) templateFunctionConfig(key string, def ...interface{}) interface{} {
	var defValue interface{}

	if len(def) > 0 {
		defValue = def[0]
	}

	if c.Has(key) {
		return c.Get(key)
	}

	return defValue
}
