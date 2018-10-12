package internal

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/internal/handlers"
)

func (c *Component) getServeMux() (*Router, error) {
	// init routes
	components, err := c.application.GetComponents()
	if err != nil {
		return nil, err
	}

	router := NewRouter(c.logger, c.config.Int(dashboard.ConfigPanicHandlerCallerSkip))

	// Special pages
	router.SetPanicHandler(&handlers.PanicHandler{})
	router.SetNotFoundHandler(&handlers.NotFoundHandler{})
	router.SetNotAllowedHandler(&handlers.MethodNotAllowedHandler{})

	// Middleware
	router.addMiddleware(ContextMiddleware(c.application, router, c.config, c.logger, c.renderer, c.session))

	for _, component := range components {
		if componentRoute, ok := component.(dashboard.HasRoutes); ok {
			for _, route := range componentRoute.DashboardRoutes() {
				router.addRoute(NewRouteItem(route, component))
			}
		}

		if componentMiddleware, ok := component.(dashboard.HasServerMiddleware); ok {
			for _, middleware := range componentMiddleware.DashboardMiddleware() {
				router.addMiddleware(middleware)
			}
		}
	}

	// fixing middleware
	router.addMiddleware(AuthorizationMiddleware)

	// fixing routes
	startUrlRoute := NewRouteItem(dashboard.NewRoute("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, c.config.String(dashboard.ConfigStartURL), http.StatusMovedPermanently)
	})), c)

	router.addRoute(startUrlRoute)

	return router, nil
}
