package profiling

import (
	"net/http"
	pprofHandlers "net/http/pprof"
	"runtime/pprof"

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
		Name: "Profiling",
		Icon: "terminal",
		SubMenu: []*dashboard.Menu{
			{
				Name: "Traces",
				Url:  "/trace",
			},
			{
				Name:   "Pprof",
				Url:    "/debug/pprof",
				Direct: true,
			},
		},
	}
}

func (c *Component) GetDashboardRoutes() []*dashboard.Route {
	routes := []*dashboard.Route{
		{
			Methods: []string{http.MethodGet, http.MethodPost},
			Path:    "/trace",
			Handler: &TraceHandler{
				config: c.config,
			},
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/debug/pprof/cmdline",
			Handler: c.debugHandler(pprofHandlers.Cmdline),
			Direct:  true,
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/debug/pprof/profile",
			Handler: c.debugHandler(pprofHandlers.Profile),
			Direct:  true,
		},
		{
			Methods: []string{http.MethodGet, http.MethodPost},
			Path:    "/debug/pprof/symbol",
			Handler: c.debugHandler(pprofHandlers.Symbol),
			Direct:  true,
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/debug/pprof/trace",
			Handler: c.debugHandler(pprofHandlers.Trace),
			Direct:  true,
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/debug/pprof/",
			Handler: c.debugHandler(pprofHandlers.Index),
			Direct:  true,
		},
	}

	for _, profile := range pprof.Profiles() {
		routes = append(routes, &dashboard.Route{
			Methods: []string{http.MethodGet},
			Path:    "/debug/pprof/" + profile.Name(),
			Handler: c.debugHandler(pprofHandlers.Index),
			Direct:  true,
		})
	}

	return routes
}
