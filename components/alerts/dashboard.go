package alerts

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
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

func (c *Component) GetDashboardMenu() *dashboard.Menu {
	return &dashboard.Menu{
		Name: "Alerts",
		Url:  "/",
		Icon: "comments",
	}
}

func (c *Component) GetDashboardRoutes() []*dashboard.Route {
	return []*dashboard.Route{
		{
			Methods: []string{http.MethodGet},
			Path:    "/",
			Handler: &IndexHandler{
				component: c,
			},
			Auth: true,
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/ajax/",
			Handler: &AjaxHandler{
				component: c,
			},
			Auth: true,
		},
	}
}
