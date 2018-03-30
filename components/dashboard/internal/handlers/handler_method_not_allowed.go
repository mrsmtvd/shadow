package handlers

import (
	"context"
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
)

type MethodNotAllowedHandler struct {
	dashboard.Handler
}

func (h *MethodNotAllowedHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)

	// FIXME: refactoring
	ctx := context.WithValue(r.Context(), dashboard.ComponentContextKey, r.Application().GetComponent(dashboard.ComponentName))

	h.RenderLayout(ctx, "405", "simple", nil)
}
