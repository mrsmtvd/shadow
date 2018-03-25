package internal

import (
	"net/http"
	pprofHandlers "net/http/pprof"
	"runtime/pprof"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/profiling"
	"github.com/kihamo/shadow/components/profiling/internal/handlers"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "templates",
	}
}

func (c *Component) DashboardMenu() dashboard.Menu {
	routes := c.DashboardRoutes()

	show := func(r *dashboard.Request) bool {
		return r.Config().Bool(config.ConfigDebug)
	}

	return dashboard.NewMenuWithUrl(
		"Profiling",
		"",
		profiling.ComponentName,
		"terminal",
		[]dashboard.Menu{
			dashboard.NewMenuWithRoute("Trace", routes[1], "", nil, show),
			dashboard.NewMenuWithRoute("Pprof", routes[7], "", nil, show),
			dashboard.NewMenuWithRoute("Expvar", routes[2], "", nil, show),
		},
		show,
	)
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.NewRoute(
				c.Name(),
				[]string{http.MethodGet},
				"/"+c.Name()+"/assets/*filepath",
				&assetfs.AssetFS{
					Asset:     Asset,
					AssetDir:  AssetDir,
					AssetInfo: AssetInfo,
					Prefix:    "assets",
				},
				"",
				false),
			dashboard.NewRoute(
				c.Name(),
				[]string{http.MethodGet, http.MethodPost},
				"/"+c.Name()+"/trace/",
				&handlers.TraceHandler{},
				"",
				true),
			dashboard.NewRoute(
				c.Name(),
				[]string{http.MethodGet},
				"/debug/vars/",
				&handlers.ExpvarHandler{},
				"",
				false),
			dashboard.NewRoute(
				c.Name(),
				[]string{http.MethodGet},
				"/debug/pprof/cmdline",
				&handlers.DebugHandler{
					HandlerFunc: pprofHandlers.Cmdline,
				},
				"",
				false),
			dashboard.NewRoute(
				c.Name(),
				[]string{http.MethodGet},
				"/debug/pprof/profile",
				&handlers.DebugHandler{
					HandlerFunc: pprofHandlers.Profile,
				},
				"",
				false),
			dashboard.NewRoute(
				c.Name(),
				[]string{http.MethodGet},
				"/debug/pprof/symbol",
				&handlers.DebugHandler{
					HandlerFunc: pprofHandlers.Symbol,
				},
				"",
				false),
			dashboard.NewRoute(
				c.Name(),
				[]string{http.MethodGet},
				"/debug/pprof/trace",
				&handlers.DebugHandler{
					HandlerFunc: pprofHandlers.Trace,
				},
				"",
				false),
			dashboard.NewRoute(
				c.Name(),
				[]string{http.MethodGet},
				"/debug/pprof/",
				&handlers.DebugHandler{
					HandlerFunc: pprofHandlers.Index,
				},
				"",
				false),
		}

		for _, profile := range pprof.Profiles() {
			c.routes = append(c.routes, dashboard.NewRoute(
				c.Name(),
				[]string{http.MethodGet},
				"/debug/pprof/"+profile.Name(),
				&handlers.DebugHandler{
					HandlerFunc: pprofHandlers.Index,
				},
				"",
				false))
		}
	}

	return c.routes
}
