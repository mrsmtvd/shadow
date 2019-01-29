package handlers

import (
	"context"
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
)

type NotFoundHandler struct {
	dashboard.Handler

	component dashboard.Component
}

func NewNotFoundHandler(component dashboard.Component) *NotFoundHandler {
	return &NotFoundHandler{
		component: component,
	}
}

func (h *NotFoundHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	w.WriteHeader(http.StatusNotFound)

	// FIXME: refactoring
	ctx := context.WithValue(r.Context(), dashboard.ComponentContextKey, h.component)

	h.RenderLayout(ctx, "404", "simple", nil)
}
