package handlers

import (
	"context"
	"time"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/database"
)

type StatusHandler struct {
	dashboard.Handler

	Component database.Component
}

func (h *StatusHandler) status(e database.Executor, ctx context.Context) string {
	c, _ := context.WithTimeout(ctx, time.Second)

	if err := e.Ping(c); err != nil {
		return err.Error()
	}

	return "OK"
}

func (h *StatusHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	s := h.Component.Storage()

	master := s.Master()
	slaves := s.Slaves()

	statuses := make([][]string, 0, 1+len(slaves))
	statuses = append(statuses, []string{master.String(), "Master", h.status(master, r.Context())})

	for _, slave := range slaves {
		statuses = append(statuses, []string{slave.String(), "Slave", h.status(slave, r.Context())})
	}

	h.Render(r.Context(), h.Component.GetName(), "status", map[string]interface{}{
		"statuses": statuses,
	})
}
