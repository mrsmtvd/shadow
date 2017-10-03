package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/internal/handlers"
)

func (c *Component) GetTemplates() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "templates",
	}
}

func (c *Component) GetDashboardMenu() dashboard.Menu {
	routes := c.GetDashboardRoutes()

	return dashboard.NewMenuItemWithRoute(
		"Dashboard",
		routes[9],
		"dashboard",
		[]dashboard.Menu{
			dashboard.NewMenuItemWithRoute("Components", routes[8], "", nil, nil),
			dashboard.NewMenuItemWithRoute("Environment", routes[2], "", nil, nil),
			dashboard.NewMenuItemWithRoute("Bindata", routes[1], "", nil, nil),
			dashboard.NewMenuItemWithRoute("Routing", routes[3], "", nil, nil),
		},
		nil,
	)
}

func (c *Component) GetDashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.NewRouteItem(
				c.GetName(),
				[]string{http.MethodGet},
				"/"+c.GetName()+"/assets/*filepath",
				&assetfs.AssetFS{
					Asset:     Asset,
					AssetDir:  AssetDir,
					AssetInfo: AssetInfo,
					Prefix:    "assets",
				},
				"",
				false),
			dashboard.NewRouteItem(
				c.GetName(),
				[]string{http.MethodGet},
				"/"+c.GetName()+"/bindata",
				&handlers.BindataHandler{
					Application: c.application,
				},
				"",
				true),
			dashboard.NewRouteItem(
				c.GetName(),
				[]string{http.MethodGet},
				"/"+c.GetName()+"/environment",
				&handlers.EnvironmentHandler{},
				"",
				true),
			dashboard.NewRouteItem(
				c.GetName(),
				[]string{http.MethodGet},
				"/"+c.GetName()+"/routing",
				&handlers.RoutingHandler{},
				"",
				true),
			dashboard.NewRouteItem(
				c.GetName(),
				[]string{http.MethodGet, http.MethodPost},
				dashboard.AuthPath+"/:provider/callback",
				&handlers.AuthHandler{
					IsCallback: true,
				},
				"",
				false),
			dashboard.NewRouteItem(
				c.GetName(),
				[]string{http.MethodGet, http.MethodPost},
				dashboard.AuthPath+"/:provider",
				&handlers.AuthHandler{},
				"",
				false),
			dashboard.NewRouteItem(
				c.GetName(),
				[]string{http.MethodGet},
				dashboard.AuthPath,
				&handlers.AuthHandler{},
				"",
				false),
			dashboard.NewRouteItem(
				c.GetName(),
				[]string{http.MethodGet},
				"/"+c.GetName()+"/logout",
				&handlers.LogoutHandler{},
				"",
				true),
		}

		componentsHandler := &handlers.ComponentsHandler{
			Application: c.application,
		}

		c.routes = append(c.routes, []dashboard.Route{
			dashboard.NewRouteItem(
				c.GetName(),
				[]string{http.MethodGet},
				"/"+c.GetName()+"/components",
				componentsHandler,
				"",
				true),
			dashboard.NewRouteItem(
				c.GetName(),
				[]string{http.MethodGet},
				"/"+c.GetName()+"/",
				componentsHandler,
				"",
				true),
		}...)
	}

	return c.routes
}
