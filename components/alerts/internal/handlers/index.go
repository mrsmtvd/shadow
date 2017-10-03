package handlers

import (
	"github.com/kihamo/shadow/components/alerts"
	"github.com/kihamo/shadow/components/dashboard"
)

type IndexHandler struct {
	dashboard.Handler

	Component alerts.Component
}

func (h *IndexHandler) ServeHTTP(_ *dashboard.Response, r *dashboard.Request) {
	h.Render(r.Context(), h.Component.GetName(), "index", map[string]interface{}{
		"alerts": h.Component.GetAlerts(),
	})
}
