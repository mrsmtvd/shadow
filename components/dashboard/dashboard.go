package dashboard

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
)

func (c *Component) GetTemplates() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "templates",
	}
}

func (c *Component) GetDashboardMenu() *Menu {
	return &Menu{
		Name: "Dashboard",
		Icon: "dashboard",
		SubMenu: []*Menu{
			{
				Name: "Components",
				Url:  "/components",
			},
			{
				Name: "Configuration",
				Url:  "/config",
			},
			{
				Name: "Environment",
				Url:  "/environment",
			},
			{
				Name: "Bindata",
				Url:  "/bindata",
			},
		},
	}
}

func (c *Component) GetDashboardRoutes() []*Route {
	routes := []*Route{
		{
			Methods: []string{http.MethodGet},
			Path:    "/assets/*filepath",
			Handler: &assetfs.AssetFS{
				Asset:     Asset,
				AssetDir:  AssetDir,
				AssetInfo: AssetInfo,
				Prefix:    "assets",
			},
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/bindata",
			Handler: &BindataHandler{
				application: c.application,
			},
		},
		{
			Methods: []string{http.MethodGet, http.MethodPost},
			Path:    "/config",
			Handler: &ConfigHandler{
				application: c.application,
			},
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/environment",
			Handler: &EnvironmentHandler{},
		},
	}

	componentsHandler := &ComponentsHandler{
		application: c.application,
	}

	routes = append(routes, []*Route{
		{
			Methods: []string{http.MethodGet},
			Path:    "/components",
			Handler: componentsHandler,
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/",
			Handler: componentsHandler,
		},
	}...)

	return routes
}
