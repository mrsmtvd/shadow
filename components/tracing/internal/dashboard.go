package internal

import (
	"net/http"

	m "github.com/kihamo/shadow/components/tracing/http"
)

func (c *Component) DashboardMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		c.serverMiddleware,
	}
}

func (c *Component) serverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := c.Tracer()

		if tracer == nil {
			next.ServeHTTP(w, r)
		} else {
			m.ServerMiddleware(tracer)(next).ServeHTTP(w, r)
		}
	})
}
