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
	return dashboard.NewMenu("Database").
		WithUrl("/" + c.Name() + "/").
		WithIcon("database").
		WithChild(dashboard.NewMenu("Migrations").WithUrl("/" + c.Name() + "/migrations/")).
		WithChild(dashboard.NewMenu("Tables").WithUrl("/" + c.Name() + "/tables/")).
		WithChild(dashboard.NewMenu("Status").WithUrl("/" + c.Name() + "/status/"))
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	return []dashboard.Route{
		dashboard.RouteFromAssetFS(c),
		dashboard.NewRoute("/"+c.Name()+"/", &handlers.StatusHandler{}).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
		dashboard.NewRoute("/"+c.Name()+"/migrations/", handlers.NewMigrationsHandler(c)).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
		dashboard.NewRoute("/"+c.Name()+"/tables/", handlers.NewTablesHandler(c)).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
		dashboard.NewRoute("/"+c.Name()+"/status/", handlers.NewStatusHandler(c)).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
	}
}
