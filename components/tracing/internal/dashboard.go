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
		m.ServerMiddleware(c.Tracer())(next).ServeHTTP(w, r)
	})
}
