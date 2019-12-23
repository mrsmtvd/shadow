package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/ota"
	"github.com/kihamo/shadow/components/ota/internal/handlers"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	routes := c.DashboardRoutes()

	return dashboard.NewMenu("OTA").
		WithIcon("cloud-upload-alt").
		WithChild(dashboard.NewMenu("Upgrade").WithRoute(routes[0])).
		WithChild(dashboard.NewMenu("Releases").WithRoute(routes[2])).
		WithChild(dashboard.NewMenu("Repository").WithRoute(routes[4]).WithShow(func(r *dashboard.Request) bool {
			return r.Config().Bool(ota.ConfigRepositoryServerEnabled)
		}))
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		releasesHandler := &handlers.ReleasesHandler{
			Updater:           c.updater,
			UploadRepository:  c.uploadRepository,
			UpgradeRepository: c.upgradeRepository,
			CurrentRelease:    c.currentRelease,
		}
		upgradeHandler := &handlers.UpgradeHandler{
			Updater:          c.updater,
			UploadRepository: c.uploadRepository,
		}
		repositoryHandler := &handlers.RepositoryHandler{
			UpgradeRepository: c.upgradeRepository,
		}

		c.routes = []dashboard.Route{
			dashboard.NewRoute("/"+c.Name()+"/upgrade/", upgradeHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/upgrade/:id/:action", upgradeHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/", releasesHandler).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/release/:id/:action", releasesHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/repository/", repositoryHandler).
				WithMethods([]string{http.MethodGet}),
			dashboard.NewRoute("/"+c.Name()+"/repository/:id/:file", repositoryHandler).
				WithMethods([]string{http.MethodGet}),
		}
	}

	return c.routes
}
