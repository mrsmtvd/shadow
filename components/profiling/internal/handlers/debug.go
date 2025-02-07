package handlers

import (
	"net/http"

	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/dashboard"
)

type DebugHandler struct {
	dashboard.Handler

	HandlerFunc http.HandlerFunc
}

func (h *DebugHandler) ServeHTTP(w http.ResponseWriter, r *dashboard.Request) {
	if !r.Config().Bool(config.ConfigDebug) {
		h.NotFound(w, r)
		return
	}

	h.HandlerFunc.ServeHTTP(w, r.Original())
}
