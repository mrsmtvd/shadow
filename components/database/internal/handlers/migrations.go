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
	migrations, err := h.Component.FindMigrations()

	h.Render(r.Context(), h.Component.GetName(), "migrations", map[string]interface{}{
		"error":      err,
		"migrations": migrations,
	})
}
