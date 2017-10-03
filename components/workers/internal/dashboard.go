package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/workers/internal/handlers"
)

func (c *Component) GetTemplates() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "templates",
	}
}

func (c *Component) GetDashboardMenu() dashboard.Menu {
	routes := c.GetDashboardRoutes()

	return dashboard.NewMenuWithRoute("Workers", routes[1], "tasks", nil, nil)
}

func (c *Component) GetDashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.NewRoute(
				c.GetName(),
				[]string{http.MethodGet},
				"/"+c.GetName()+"/assets/*filepath",
				&assetfs.AssetFS{
					Asset:     Asset,
					AssetDir:  AssetDir,
					AssetInfo: AssetInfo,
					Prefix:    "assets",
				},
				"",
				false),
			dashboard.NewRoute(
				c.GetName(),
				[]string{http.MethodGet, http.MethodPost},
				"/"+c.GetName()+"/",
				&handlers.ManagerHandler{
					Component: c,
				},
				"",
				true),
		}
	}

	return c.routes
}
