package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/alerts/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
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

	return dashboard.NewMenuItemWithRoute("Alerts", routes[0], "comments", nil, nil)
}

func (c *Component) GetDashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.NewRouteItem(
				c.GetName(),
				[]string{http.MethodGet},
				"/"+c.GetName()+"/",
				&handlers.IndexHandler{
					Component: c,
				},
				"",
				true),
			dashboard.NewRouteItem(
				c.GetName(),
				[]string{http.MethodGet},
				"/"+c.GetName()+"/ajax/",
				&handlers.AjaxHandler{
					Component: c,
				},
				"",
				true),
		}
	}

	return c.routes
}
