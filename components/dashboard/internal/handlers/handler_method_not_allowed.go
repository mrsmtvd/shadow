package handlers

import (
	"net/http"

	"github.com/mrsmtvd/shadow/components/dashboard"
)

type MethodNotAllowedHandler struct {
	dashboard.Handler
}

func (h *MethodNotAllowedHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)

	ctx := dashboard.ContextWithTemplateNamespace(r.Context(), dashboard.ComponentName)
	h.RenderLayout(ctx, "405", "simple", nil)
}
