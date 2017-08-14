package dashboard

import (
	"net/http"
)

type MethodNotAllowedHandler struct {
	Handler
}

func (h *MethodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	h.RenderLayout(r.Context(), ComponentName, "405", "simple", nil)
}
