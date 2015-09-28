package frontend

import (
	"net/http"
)

type MethodNotAllowedHandler struct {
	AbstractFrontendHandler
}

func (h *MethodNotAllowedHandler) Handle() {
	h.SetTemplate("405.tpl.html")
	h.SetPageTitle("Method not allowed")

	h.Output.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.Output.WriteHeader(http.StatusMethodNotAllowed)
}
