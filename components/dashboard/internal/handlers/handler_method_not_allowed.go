package handlers

import (
	"context"
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
)

type MethodNotAllowedHandler struct {
	dashboard.Handler

	component dashboard.Component
}

func NewMethodNotAllowedHandler(component dashboard.Component) *MethodNotAllowedHandler {
	return &MethodNotAllowedHandler{
		component: component,
	}
}

func (h *MethodNotAllowedHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)

	// FIXME: refactoring
	ctx := context.WithValue(r.Context(), dashboard.ComponentContextKey, h.component)

	h.RenderLayout(ctx, "405", "simple", nil)
}
