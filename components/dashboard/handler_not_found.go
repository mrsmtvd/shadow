package dashboard

import (
	"net/http"
)

type NotFoundHandler struct {
	Handler
}

func (h *NotFoundHandler) ServeHTTP(w *Response, r *Request) {
	w.WriteHeader(http.StatusNotFound)
	h.RenderLayout(r.Context(), ComponentName, "404", "simple", nil)
}
