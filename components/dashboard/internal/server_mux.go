package internal

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/internal/handlers"
	"github.com/kihamo/shadow/components/metrics"
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
	router.AddMiddleware(ContextMiddleware(router, c.config, c.logger, c.renderer, c.session))

	if c.application.HasComponent(metrics.ComponentName) {
		router.AddMiddleware(MetricsMiddleware())
	}

	router.AddMiddleware(LoggerMiddleware())
	router.AddMiddleware(AuthorizationMiddleware())

	for _, component := range components {
		if componentRoute, ok := component.(dashboard.HasRoutes); ok {
			for _, route := range componentRoute.DashboardRoutes() {
				router.AddRoute(route, component.Name())
			}
		}
	}

	router.AddRoute(dashboard.NewRoute("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, c.config.String(dashboard.ConfigStartURL), http.StatusMovedPermanently)
	})), c.Name())

	return router, nil
}
