package profiling

import (
	"net/http"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

type DebugHandler struct {
	dashboard.Handler

	handler http.HandlerFunc
}

func (h *DebugHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if !r.Config().GetBool(config.ConfigDebug) {
		h.NotFound(w, r)
		return
	}

	h.handler.ServeHTTP(w, r.Original())
}
