package handlers

import (
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
		"stack": error.Stack,
		"file":  error.File,
		"line":  error.Line,
	}

	w.WriteHeader(http.StatusInternalServerError)
	h.RenderLayout(r.Context(), dashboard.ComponentName, "500", "simple", map[string]interface{}{
		"panic": fields,
	})

	r.Logger().Error("Frontend request error", fields)
}
