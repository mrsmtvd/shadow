package handlers

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/auth"
)

type LogoutHandler struct {
	dashboard.Handler
}

func (h *LogoutHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	session := r.Session()
	session.Remove(dashboard.SessionUser)

	for _, provider := range auth.GetProviders() {
		session.Remove(dashboard.SessionAuthProvider(provider))
	}

	h.Redirect("/", http.StatusFound, w, r)
}
