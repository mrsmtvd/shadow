package internal

import (
	"net/http"

	m "github.com/mrsmtvd/shadow/components/tracing/http"
)

func (c *Component) DashboardMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		m.ServerMiddleware(c.Tracer()),
	}
}
