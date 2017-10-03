package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/metrics/internal/handlers"
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

	return dashboard.NewMenuItemWithRoute("Metrics", routes[0], "thermometer-empty", nil, nil)
}

func (c *Component) GetDashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.NewRouteItem(
				c.GetName(),
				[]string{http.MethodGet},
				"/"+c.GetName()+"/",
				&handlers.IndextHandler{
					Component: c,
				},
				"",
				true),
		}
	}

	return c.routes
}
