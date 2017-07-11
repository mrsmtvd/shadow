package workers

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
		Name: "Workers",
		Url:  "/",
		Icon: "tasks",
	}
}

func (c *Component) GetDashboardRoutes() []*dashboard.Route {
	return []*dashboard.Route{
		{
			Methods: []string{http.MethodGet},
			Path:    "/js/*filepath",
			Handler: &assetfs.AssetFS{
				Asset:     Asset,
				AssetDir:  AssetDir,
				AssetInfo: AssetInfo,
				Prefix:    "public/js",
			},
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/",
			Handler: &IndexHandler{},
		},
		{
			Methods: []string{http.MethodGet, http.MethodPost},
			Path:    "/ajax/",
			Handler: &AjaxHandler{
				component: c,
			},
		},
	}
}
