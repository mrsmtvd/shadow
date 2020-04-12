package handlers

import (
	"net/http"
	"runtime/debug"

	"github.com/kihamo/shadow/components/dashboard"
)

type DependenciesHandler struct {
	dashboard.Handler
}

func (h *DependenciesHandler) ServeHTTP(w http.ResponseWriter, r *dashboard.Request) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		h.NotFound(w, r)
	}

	h.Render(r.Context(), "dependencies", map[string]interface{}{
		"dependencies": info,
	})
}
