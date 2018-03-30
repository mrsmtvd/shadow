package handlers

import (
	"context"
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
)

type NotFoundHandler struct {
	dashboard.Handler
}

func (h *NotFoundHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	w.WriteHeader(http.StatusNotFound)

	// FIXME: refactoring
	ctx := context.WithValue(r.Context(), dashboard.ComponentContextKey, r.Application().GetComponent(dashboard.ComponentName))

	h.RenderLayout(ctx, "404", "simple", nil)
}
