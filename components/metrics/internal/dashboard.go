package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/metrics/internal/handlers"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	routes := c.DashboardRoutes()

	return dashboard.NewMenu("Metrics").WithRoute(routes[0]).WithIcon("thermometer-empty")
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.NewRoute("/"+c.Name()+"/list/", &handlers.ListHandler{
				Component: c,
			}).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
		}
	}

	return c.routes
}
