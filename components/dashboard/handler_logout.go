package dashboard

import (
	"net/http"

	"github.com/alexedwards/scs/session"
)

type LogoutHandler struct {
	Handler
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session.Remove(r, "username")
	h.Redirect("/", http.StatusFound, w, r)
	return
}
