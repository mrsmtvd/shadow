package internal

import (
	"net/http"
	pprofHandlers "net/http/pprof"
	"runtime/pprof"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
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

	return dashboard.NewMenu("Profiling").
		WithIcon("terminal").
		WithShow(show).
		WithChild(dashboard.NewMenu("Trace").WithRoute(routes[1])).
		WithChild(dashboard.NewMenu("Pprof").WithRoute(routes[7])).
		WithChild(dashboard.NewMenu("Expvar").WithRoute(routes[2]))
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.NewRoute(
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
				[]string{http.MethodGet, http.MethodPost},
				"/"+c.Name()+"/trace/",
				&handlers.TraceHandler{},
				"",
				true),
			dashboard.NewRoute(
				[]string{http.MethodGet},
				"/debug/vars/",
				&handlers.ExpvarHandler{},
				"",
				false),
			dashboard.NewRoute(
				[]string{http.MethodGet},
				"/debug/pprof/cmdline",
				&handlers.DebugHandler{
					HandlerFunc: pprofHandlers.Cmdline,
				},
				"",
				false),
			dashboard.NewRoute(
				[]string{http.MethodGet},
				"/debug/pprof/profile",
				&handlers.DebugHandler{
					HandlerFunc: pprofHandlers.Profile,
				},
				"",
				false),
			dashboard.NewRoute(
				[]string{http.MethodGet},
				"/debug/pprof/symbol",
				&handlers.DebugHandler{
					HandlerFunc: pprofHandlers.Symbol,
				},
				"",
				false),
			dashboard.NewRoute(
				[]string{http.MethodGet},
				"/debug/pprof/trace",
				&handlers.DebugHandler{
					HandlerFunc: pprofHandlers.Trace,
				},
				"",
				false),
			dashboard.NewRoute(
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
