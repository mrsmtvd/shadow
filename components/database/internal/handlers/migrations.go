package handlers

import (
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/database"
	"github.com/kihamo/shadow/components/i18n"
)

// easyjson:json
type migrationsHandlerResponse struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

type MigrationsManager interface {
	UpMigration(id, source string) error
	UpMigrations() (n int, err error)
	DownMigration(id, source string) error
	DownMigrations() (n int, err error)
}

type MigrationsHandler struct {
	dashboard.Handler

	component database.Component
	manager   MigrationsManager
}

func NewMigrationsHandler(component database.Component, manager MigrationsManager) *MigrationsHandler {
	return &MigrationsHandler{
		component: component,
		manager:   manager,
	}
}

func (h *MigrationsHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if r.IsPost() {
		var err error

		q := r.URL().Query()

		action := q.Get(":action")
		source := q.Get(":source")
		id := q.Get(":id")

		if action != "up" && action != "down" {
			h.NotFound(w, r)
			return
		}

		locale := i18n.Locale(r.Context())

		switch action {
		case "up":
			if source == "" && id == "" {
				_, err = h.manager.UpMigrations()
				if err == nil {
					w.SendJSON(migrationsHandlerResponse{
						Result:  "success",
						Message: locale.Translate(h.component.Name(), "Apply migrations success", ""),
					})
					return
				}
			} else if source != "" && id != "" {
				err = h.manager.UpMigration(id, source)
				if err == nil {
					w.SendJSON(migrationsHandlerResponse{
						Result:  "success",
						Message: locale.Translate(h.component.Name(), "Apply migration %s for %s success", "", id, source),
					})
					return
				}
			} else {
				h.NotFound(w, r)
				return
			}

		case "down":
			if source == "" && id == "" {
				_, err = h.manager.DownMigrations()
				if err == nil {
					w.SendJSON(migrationsHandlerResponse{
						Result:  "success",
						Message: locale.Translate(h.component.Name(), "Rollback migrations success", ""),
					})
					return
				}
			} else if source != "" && id != "" {
				err = h.manager.DownMigration(id, source)
				if err == nil {
					w.SendJSON(migrationsHandlerResponse{
						Result:  "success",
						Message: locale.Translate(h.component.Name(), "Rollback migration %s for %s success", "", id, source),
					})
					return
				}
			} else {
				h.NotFound(w, r)
				return
			}

		default:
			h.NotFound(w, r)
			return
		}

		if err != nil {
			w.SendJSON(migrationsHandlerResponse{
				Result:  "failed",
				Message: err.Error(),
			})
			return
		}
	}

	h.Render(r.Context(), "migrations", map[string]interface{}{
		"migrations": h.component.Migrations(),
	})
}
