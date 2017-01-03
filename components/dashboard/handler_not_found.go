package dashboard

import "net/http"

type NotFoundHandler struct {
	TemplateHandler
}

func (h *NotFoundHandler) Handle() {
	h.SetView("dashboard", "404")

	h.Response().WriteHeader(http.StatusNotFound)
}
