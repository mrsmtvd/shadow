package internal

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/internal/handlers"
)

func (c *Component) initServeMux() error {
	// Special pages
	c.router.SetPanicHandler(&handlers.PanicHandler{})
	c.router.SetNotFoundHandler(&handlers.NotFoundHandler{})
	c.router.SetNotAllowedHandler(&handlers.MethodNotAllowedHandler{})

	// Middleware
	c.router.addMiddleware(SessionMiddleware(c.sessionManager))
	c.router.addMiddleware(ContextMiddleware(c.router, c.renderer))

	for _, component := range c.components {
		if componentRoute, ok := component.(dashboard.HasRoutes); ok {
			for _, route := range componentRoute.DashboardRoutes() {
				c.router.addRoute(NewRouteItem(route, component))
			}
		}

		if componentMiddleware, ok := component.(dashboard.HasServerMiddleware); ok {
			for _, middleware := range componentMiddleware.DashboardMiddleware() {
				c.router.addMiddleware(middleware)
			}
		}
	}

	// fixing middleware
	c.router.addMiddleware(AuthorizationMiddleware)

	// fixing routes
	startURLRoute := NewRouteItem(dashboard.NewRoute("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, c.config.String(dashboard.ConfigStartURL), http.StatusMovedPermanently)
	})), c)

	c.router.addRoute(startURLRoute)

	return nil
}
