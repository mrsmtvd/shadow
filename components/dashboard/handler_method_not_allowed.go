package dashboard

import (
	"net/http"
)

type MethodNotAllowedHandler struct {
	Handler
}

func (h *MethodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	h.Render(r.Context(), "dashboard", "405", nil)
}
