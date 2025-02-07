package handlers

import (
	"net/http"

	"github.com/mrsmtvd/shadow/components/dashboard"
)

type NotFoundHandler struct {
	dashboard.Handler
}

func (h *NotFoundHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	w.WriteHeader(http.StatusNotFound)

	ctx := dashboard.ContextWithTemplateNamespace(r.Context(), dashboard.ComponentName)
	h.RenderLayout(ctx, "404", "simple", nil)
}
