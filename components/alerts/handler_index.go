package alerts

import (
	"github.com/kihamo/shadow/components/dashboard"
)

type IndexHandler struct {
	dashboard.Handler

	component *Component
}

func (h *IndexHandler) ServeHTTP(_ *dashboard.Response, r *dashboard.Request) {
	h.Render(r.Context(), ComponentName, "index", map[string]interface{}{
		"alerts": h.component.GetAlerts(),
	})
}
