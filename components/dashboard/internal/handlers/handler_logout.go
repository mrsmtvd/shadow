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

	if err := session.Remove(dashboard.SessionUser); err != nil {
		panic(err.Error())
	}

	if err := session.Remove(dashboard.AuthSessionName()); err != nil {
		panic(err.Error())
	}

	h.Redirect(r.Config().String(dashboard.ConfigStartURL), http.StatusFound, w, r)
}
