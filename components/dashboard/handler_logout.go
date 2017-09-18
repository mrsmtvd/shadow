package dashboard

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard/auth"
)

type LogoutHandler struct {
	Handler
}

func (h *LogoutHandler) ServeHTTP(w *Response, r *Request) {
	session := r.Session()
	session.Remove(SessionUser)

	for _, provider := range auth.GetProviders() {
		session.Remove(SessionAuthProvider(provider))
	}

	h.Redirect("/", http.StatusFound, w, r)
}
