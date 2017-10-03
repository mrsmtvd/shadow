package internal

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/internal/handlers"
	"github.com/kihamo/shadow/components/metrics"
)

func (c *Component) getServeMux() (http.Handler, error) {
	// init routes
	components, err := c.application.GetComponents()
	if err != nil {
		return nil, err
	}

	router := NewRouter(c)

	// Special pages
	router.SetPanicHandler(&handlers.PanicHandler{})
	router.SetNotFoundHandler(&handlers.NotFoundHandler{})
	router.SetNotAllowedHandler(&handlers.MethodNotAllowedHandler{})

	// Middleware
	router.AddMiddleware(ContextMiddleware(router, c.config, c.logger, c.renderer, c.session))

	if c.application.HasComponent(metrics.ComponentName) {
		router.AddMiddleware(MetricsMiddleware())
	}

	router.AddMiddleware(LoggerMiddleware())
	router.AddMiddleware(AuthorizationMiddleware())

	for _, component := range components {
		if componentRoute, ok := component.(dashboard.HasRoutes); ok {
			for _, route := range componentRoute.GetDashboardRoutes() {
				router.AddRoute(route)
			}
		}
	}

	router.AddRoute(dashboard.NewRouteItem(
		dashboard.ComponentName,
		nil,
		"/",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/"+c.GetName(), http.StatusMovedPermanently)
		}),
		"",
		false))

	return router, nil
}
