package handlers

import (
	"net/http"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

type DebugHandler struct {
	dashboard.Handler

	HandlerFunc http.HandlerFunc
}

func (h *DebugHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if !r.Config().Bool(config.ConfigDebug) {
		h.NotFound(w, r)
		return
	}

	h.HandlerFunc.ServeHTTP(w, r.Original())
}
