package internal

import (
	"net/http"
	pprofHandlers "net/http/pprof"
	"runtime/pprof"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/profiling/internal/handlers"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	show := func(r *dashboard.Request) bool {
		return r.Config().Bool(config.ConfigDebug)
	}

	return dashboard.NewMenu("Profiling").
		WithIcon("terminal").
		WithShow(show).
		WithChild(dashboard.NewMenu("Trace").WithUrl("/" + c.Name() + "/trace/")).
		WithChild(dashboard.NewMenu("Pprof").WithUrl("/debug/pprof/")).
		WithChild(dashboard.NewMenu("Expvar").WithUrl("/debug/vars/"))
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	routes := []dashboard.Route{
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
		routes = append(routes, dashboard.NewRoute("/debug/pprof/"+profile.Name(), &handlers.DebugHandler{
			HandlerFunc: pprofHandlers.Index,
		}).
			WithMethods([]string{http.MethodGet}))
	}

	return routes
}
