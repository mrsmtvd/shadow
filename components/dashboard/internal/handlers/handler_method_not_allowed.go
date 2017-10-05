package handlers

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
)

type MethodNotAllowedHandler struct {
	dashboard.Handler
}

func (h *MethodNotAllowedHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	h.RenderLayout(r.Context(), dashboard.ComponentName, "405", "simple", nil)
}
