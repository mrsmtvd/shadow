package internal

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logging"
	m "github.com/kihamo/shadow/components/logging/http"
)

func (c *Component) DashboardMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// save logger in context
				if request := dashboard.RequestFromContext(r.Context()); request != nil {
					name := dashboard.TemplateNamespaceFromContext(r.Context())
					request.WithContext(logging.ContextWithLogger(r.Context(), c.Logger().Named(name)))
					r = request.Original()
				}

				next.ServeHTTP(w, r)
			})
		},
		m.ServerMiddleware(c.Logger().Named(dashboard.ComponentName)),
	}
}
