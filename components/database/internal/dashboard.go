package internal

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
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
	migrationsHandler := handlers.NewMigrationsHandler(c, c)

	return []dashboard.Route{
		dashboard.RouteFromAssetFS(c),
		dashboard.NewRoute("/"+c.Name()+"/", &handlers.StatusHandler{}).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
		dashboard.NewRoute("/"+c.Name()+"/migrations/", migrationsHandler).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
		dashboard.NewRoute("/"+c.Name()+"/migrations/:action/", migrationsHandler).
			WithMethods([]string{http.MethodPost}).
			WithAuth(true),
		dashboard.NewRoute("/"+c.Name()+"/migrations/:action/:source/:id", migrationsHandler).
			WithMethods([]string{http.MethodPost}).
			WithAuth(true),
		dashboard.NewRoute("/"+c.Name()+"/tables/", handlers.NewTablesHandler(c)).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
		dashboard.NewRoute("/"+c.Name()+"/status/", handlers.NewStatusHandler(c)).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
	}
}
