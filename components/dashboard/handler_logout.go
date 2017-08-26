package dashboard

import (
	"net/http"
)

type LogoutHandler struct {
	Handler
}

func (h *LogoutHandler) ServeHTTP(w *Response, r *Request) {
	r.Session().Remove("username")
	h.Redirect("/", http.StatusFound, w, r)
	return
}
