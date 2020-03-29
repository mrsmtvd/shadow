package internal

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
	m "github.com/kihamo/shadow/components/metrics/http"
	"github.com/kihamo/shadow/components/metrics/internal/handlers"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	return dashboard.NewMenu("Metrics").
		WithURL("/" + c.Name() + "/list/").
		WithIcon("thermometer-empty")
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	return []dashboard.Route{
		dashboard.NewRoute("/"+c.Name()+"/list/", handlers.NewListHandler(c)).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
	}
}

func (c *Component) DashboardMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		m.ServerMiddleware,
	}
}
