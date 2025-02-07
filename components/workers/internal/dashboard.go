package internal

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/workers/internal/handlers"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	return dashboard.NewMenu("Workers").
		WithURL("/" + c.Name() + "/").
		WithIcon("tasks")
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	return []dashboard.Route{
		dashboard.RouteFromAssetFS(c),
		dashboard.NewRoute("/"+c.Name()+"/", handlers.NewManagerHandler(c)).
			WithMethods([]string{http.MethodGet, http.MethodPost}).
			WithAuth(true),
	}
}
