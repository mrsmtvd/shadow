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
			{
				Name:   "Expvar",
				Url:    "/debug/vars",
				Direct: true,
			},
		},
	}
}

func (c *Component) GetDashboardRoutes() []*dashboard.Route {
	routes := []*dashboard.Route{
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
			Methods: []string{http.MethodGet, http.MethodPost},
			Path:    "/trace",
			Handler: &TraceHandler{},
			Auth:    true,
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/debug/vars/",
			Handler: &ExpvarHandler{},
			Direct:  true,
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/debug/pprof/cmdline",
			Handler: &DebugHandler{
				handler: pprofHandlers.Cmdline,
			},
			Direct: true,
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/debug/pprof/profile",
			Handler: &DebugHandler{
				handler: pprofHandlers.Profile,
			},
			Direct: true,
		},
		{
			Methods: []string{http.MethodGet, http.MethodPost},
			Path:    "/debug/pprof/symbol",
			Handler: &DebugHandler{
				handler: pprofHandlers.Symbol,
			},
			Direct: true,
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/debug/pprof/trace",
			Handler: &DebugHandler{
				handler: pprofHandlers.Trace,
			},
			Direct: true,
		},
		{
			Methods: []string{http.MethodGet},
			Path:    "/debug/pprof/",
			Handler: &DebugHandler{
				handler: pprofHandlers.Index,
			},
			Direct: true,
		},
	}

	for _, profile := range pprof.Profiles() {
		routes = append(routes, &dashboard.Route{
			Methods: []string{http.MethodGet},
			Path:    "/debug/pprof/" + profile.Name(),
			Handler: &DebugHandler{
				handler: pprofHandlers.Index,
			},
			Direct: true,
		})
	}

	return routes
}
