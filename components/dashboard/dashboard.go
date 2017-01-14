package dashboard

import (
	"net/http"
	"net/http/pprof"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/config"
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
		},
	}
}

func (c *Component) GetDashboardRoutes() []*Route {
	routes := []*Route{
		{
			Methods: []string{http.MethodGet},
			Path:    "/css/*filepath",
			Handler: &assetfs.AssetFS{
				Asset:     Asset,
				AssetDir:  AssetDir,
				AssetInfo: AssetInfo,
				Prefix:    "public/css",
			},
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/images/*filepath",
			Handler: &assetfs.AssetFS{
				Asset:     Asset,
				AssetDir:  AssetDir,
				AssetInfo: AssetInfo,
				Prefix:    "public/images",
			},
		},
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
			Path:    "/vendor/*filepath",
			Handler: &assetfs.AssetFS{
				Asset:     Asset,
				AssetDir:  AssetDir,
				AssetInfo: AssetInfo,
				Prefix:    "public/vendor",
			},
		},
		{
			Methods: []string{http.MethodGet, http.MethodPost},
			Path:    "/config",
			Handler: &ConfigHandler{},
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

	if c.config.GetBool(config.ConfigDebug) {
		routes = append(routes, []*Route{
			{
				Methods: []string{http.MethodGet},
				Path:    "/debug/pprof/cmdline",
				Handler: c.debugHandler(pprof.Cmdline),
			},
			{
				Methods: []string{http.MethodGet},
				Path:    "/debug/pprof/profile",
				Handler: c.debugHandler(pprof.Profile),
			},
			{
				Methods: []string{http.MethodGet, http.MethodPost},
				Path:    "/debug/pprof/symbol",
				Handler: c.debugHandler(pprof.Symbol),
			},
			{
				Methods: []string{http.MethodGet},
				Path:    "/debug/pprof/block",
				Handler: c.debugHandler(pprof.Index),
			},
			{
				Methods: []string{http.MethodGet},
				Path:    "/debug/pprof/goroutine",
				Handler: c.debugHandler(pprof.Index),
			},
			{
				Methods: []string{http.MethodGet},
				Path:    "/debug/pprof/heap",
				Handler: c.debugHandler(pprof.Index),
			},
			{
				Methods: []string{http.MethodGet},
				Path:    "/debug/pprof/threadcreate",
				Handler: c.debugHandler(pprof.Index),
			},
			{
				Methods: []string{http.MethodGet},
				Path:    "/debug/pprof/",
				Handler: c.debugHandler(pprof.Index),
			},
		}...)
	}

	return routes
}

func (c *Component) debugHandler(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c.config.GetBool(config.ConfigDebug) {
			h.ServeHTTP(w, r)
		} else {
			c.router.NotFound.ServeHTTP(w, r)
		}
	})
}
