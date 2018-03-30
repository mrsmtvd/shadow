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
	return dashboard.TemplatesFromAssetFS(c)
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
			dashboard.RouteFromAssetFS(c),
			dashboard.NewRoute("/"+c.Name()+"/trace/", &handlers.TraceHandler{}).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/debug/vars/", &handlers.ExpvarHandler{}).
				WithMethods([]string{http.MethodGet}),
			dashboard.NewRoute("/debug/pprof/cmdline", &handlers.DebugHandler{
				HandlerFunc: pprofHandlers.Cmdline,
			}).
				WithMethods([]string{http.MethodGet}),
			dashboard.NewRoute("/debug/pprof/profile", &handlers.DebugHandler{
				HandlerFunc: pprofHandlers.Profile,
			}).
				WithMethods([]string{http.MethodGet}),
			dashboard.NewRoute("/debug/pprof/symbol", &handlers.DebugHandler{
				HandlerFunc: pprofHandlers.Symbol,
			}).
				WithMethods([]string{http.MethodGet}),
			dashboard.NewRoute("/debug/pprof/trace", &handlers.DebugHandler{
				HandlerFunc: pprofHandlers.Trace,
			}).
				WithMethods([]string{http.MethodGet}),
			dashboard.NewRoute("/debug/pprof/", &handlers.DebugHandler{
				HandlerFunc: pprofHandlers.Index,
			}).
				WithMethods([]string{http.MethodGet}),
		}

		for _, profile := range pprof.Profiles() {
			c.routes = append(c.routes, dashboard.NewRoute("/debug/pprof/"+profile.Name(), &handlers.DebugHandler{
				HandlerFunc: pprofHandlers.Index,
			}).
				WithMethods([]string{http.MethodGet}))
		}
	}

	return c.routes
}
