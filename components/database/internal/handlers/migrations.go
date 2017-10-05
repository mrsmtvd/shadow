package handlers

import (
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/database"
)

type MigrationsHandler struct {
	dashboard.Handler

	Component database.Component
}

func (h *MigrationsHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	h.Render(r.Context(), h.Component.GetName(), "migrations", map[string]interface{}{
		"migrations": h.Component.GetAllMigrations(),
	})
}
