package internal

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	m "github.com/kihamo/shadow/components/logging/http"
)

func (c *Component) DashboardMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		m.ServerMiddleware(c.Get(dashboard.ComponentName)),
	}
}
