package handlers

import (
	"expvar"
	"net/http"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

type ExpvarHandler struct {
	dashboard.Handler
}

func (h *ExpvarHandler) ServeHTTP(w http.ResponseWriter, r *dashboard.Request) {
	if !r.Config().Bool(config.ConfigDebug) {
		h.NotFound(w, r)
		return
	}

	expvar.Handler().ServeHTTP(w, r.Original())
}
