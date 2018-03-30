package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/config/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "templates",
	}
}

func (c *Component) DashboardMenu() dashboard.Menu {
	routes := c.DashboardRoutes()

	return dashboard.NewMenu("Configuration").WithRoute(routes[1]).WithIcon("cog")
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.NewRoute(
				[]string{http.MethodGet},
				"/"+c.Name()+"/assets/*filepath",
				&assetfs.AssetFS{
					Asset:     Asset,
					AssetDir:  AssetDir,
					AssetInfo: AssetInfo,
					Prefix:    "assets",
				},
				"",
				false),
			dashboard.NewRoute(
				[]string{http.MethodGet, http.MethodPost},
				"/"+c.Name()+"/",
				&handlers.ManagerHandler{
					Application: c.application,
					Component:   c,
				},
				"",
				true),
		}
	}

	return c.routes
}

func (c *Component) DashboardTemplateFunctions() map[string]interface{} {
	return map[string]interface{}{
		"config": c.templateFunctionConfig,
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
