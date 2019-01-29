package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
)

type PanicHandler struct {
	dashboard.Handler

	component dashboard.Component
}

func NewPanicHandler(component dashboard.Component) *PanicHandler {
	return &PanicHandler{
		component: component,
	}
}

func (h *PanicHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	error := dashboard.PanicFromContext(r.Context())
	fields := map[string]interface{}{
		"error": fmt.Sprintf("%v", error.Error),
		"stack": string(error.Stack),
		"file":  error.File,
		"line":  error.Line,
	}

	w.WriteHeader(http.StatusInternalServerError)

	// FIXME: refactoring
	ctx := context.WithValue(r.Context(), dashboard.ComponentContextKey, h.component)

	h.RenderLayout(ctx, "500", "simple", map[string]interface{}{
		"panic": fields,
	})

	h.Logger().Error("Internal server error",
		"error", fields["error"],
		"stack", fields["stack"],
		"file", fields["file"],
		"line", fields["line"],
	)
}
