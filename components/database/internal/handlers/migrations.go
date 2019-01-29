package handlers

import (
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/database"
)

type MigrationsHandler struct {
	dashboard.Handler

	component database.Component
}

func NewMigrationsHandler(component database.Component) *MigrationsHandler {
	return &MigrationsHandler{
		component: component,
	}
}

func (h *MigrationsHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	h.Render(r.Context(), "migrations", map[string]interface{}{
		"migrations": h.component.Migrations(),
	})
}
