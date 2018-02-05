package handlers

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
)

type LogoutHandler struct {
	dashboard.Handler
}

func (h *LogoutHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	session := r.Session()
	session.Remove(dashboard.SessionUser)
	session.Remove(dashboard.AuthSessionName())

	h.Redirect(r.Config().String(dashboard.ConfigStartURL), http.StatusFound, w, r)
}
