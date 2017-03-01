package alerts

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
)

type IndexHandler struct {
	dashboard.Handler

	component *Component
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Render(r.Context(), ComponentName, "index", map[string]interface{}{
		"alerts": h.component.GetAlerts(),
	})
}
