package handlers

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
)

type NotFoundHandler struct {
	dashboard.Handler
}

func (h *NotFoundHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	w.WriteHeader(http.StatusNotFound)
	h.RenderLayout(r.Context(), dashboard.ComponentName, "404", "simple", nil)
}
