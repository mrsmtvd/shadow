package handlers

import (
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/database"
)

type MigrationsHandler struct {
	dashboard.Handler
}

func (h *MigrationsHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	h.Render(r.Context(), "migrations", map[string]interface{}{
		"migrations": r.Component().(database.Component).Migrations(),
	})
}
