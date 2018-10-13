package internal

import (
	"net/http"

	m "github.com/kihamo/shadow/components/tracing/http"
)

func (c *Component) DashboardMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		m.ServerMiddleware(c.Tracer()),
	}
}
