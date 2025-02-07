package internal

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/database/internal/handlers"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	return dashboard.NewMenu("Database").
		WithURL("/" + c.Name() + "/").
		WithIcon("database").
		WithChild(dashboard.NewMenu("Migrations").WithURL("/" + c.Name() + "/migrations/")).
		WithChild(dashboard.NewMenu("Tables").WithURL("/" + c.Name() + "/tables/")).
		WithChild(dashboard.NewMenu("Status").WithURL("/" + c.Name() + "/status/"))
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
