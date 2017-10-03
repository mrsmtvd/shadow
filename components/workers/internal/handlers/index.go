package handlers

import (
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/workers"
)

type IndexHandler struct {
	dashboard.Handler

	Component workers.Component
}

func (h *IndexHandler) ServeHTTP(_ *dashboard.Response, r *dashboard.Request) {
	h.Render(r.Context(), h.Component.GetName(), "index", map[string]interface{}{
		"defaultListenerName": h.Component.GetDefaultListenerName(),
	})
}
