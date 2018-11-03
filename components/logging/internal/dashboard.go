package internal

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logging"
	m "github.com/kihamo/shadow/components/logging/http"
)

func (c *Component) DashboardMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		m.ServerMiddleware(logging.DefaultLogger().Named(dashboard.ComponentName)),
	}
}
