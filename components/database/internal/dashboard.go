package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/database/internal/handlers"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	routes := c.DashboardRoutes()

	return dashboard.NewMenu("Database").
		WithRoute(routes[1]).
		WithIcon("database").
		WithChild(dashboard.NewMenu("Migrations").WithRoute(routes[2])).
		WithChild(dashboard.NewMenu("Status").WithRoute(routes[3]))
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.RouteFromAssetFS(c),
			dashboard.NewRoute("/"+c.Name(), &handlers.StatusHandler{
				Component: c,
			}).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/migrations/", &handlers.MigrationsHandler{
				Component: c,
			}).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/status/", &handlers.StatusHandler{
				Component: c,
			}).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
		}
	}

	return c.routes
}
