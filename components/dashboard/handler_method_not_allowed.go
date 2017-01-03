package dashboard

import "net/http"

type MethodNotAllowedHandler struct {
	TemplateHandler
}

func (h *MethodNotAllowedHandler) Handle() {
	h.SetView("dashboard", "405")

	h.Response().WriteHeader(http.StatusMethodNotAllowed)
}
