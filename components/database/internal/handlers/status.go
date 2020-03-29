package handlers

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/database"
)

type StatusHandler struct {
	dashboard.Handler

	component database.Component
}

func NewStatusHandler(component database.Component) *StatusHandler {
	return &StatusHandler{
		component: component,
	}
}

func (h *StatusHandler) status(ctx context.Context, e driver.Pinger) string {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := e.Ping(ctx); err != nil {
		return err.Error()
	}

	return "OK"
}

func (h *StatusHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	s := h.component.Storage()

	master := s.Master()
	slaves := s.Slaves()

	statuses := make([][]string, 0, 1+len(slaves))
	statuses = append(statuses, []string{master.String(), "Master", h.status(r.Context(), master)})

	for _, slave := range slaves {
		statuses = append(statuses, []string{slave.String(), "Slave", h.status(r.Context(), slave)})
	}

	h.Render(r.Context(), "status", map[string]interface{}{
		"statuses": statuses,
	})
}
