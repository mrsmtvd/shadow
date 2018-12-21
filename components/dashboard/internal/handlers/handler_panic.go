package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
)

type PanicHandler struct {
	dashboard.Handler
}

func (h *PanicHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	error := dashboard.PanicFromContext(r.Context())
	fields := map[string]interface{}{
		"error": fmt.Sprintf("%s", error.Error),
		"stack": string(error.Stack),
		"file":  error.File,
		"line":  error.Line,
	}

	w.WriteHeader(http.StatusInternalServerError)

	// FIXME: refactoring
	ctx := context.WithValue(r.Context(), dashboard.ComponentContextKey, r.Application().GetComponent(dashboard.ComponentName))

	h.RenderLayout(ctx, "500", "simple", map[string]interface{}{
		"panic": fields,
	})

	h.Logger().Error("Frontend request error", fields)
}
