package frontend

import (
	"net/http"
)

type NotFoundHandler struct {
	AbstractFrontendHandler
}

func (h *NotFoundHandler) Handle() {
	h.SetTemplate("404.tpl.html")
	h.View.Context["PageTitle"] = "Page not found"

	h.Output.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.Output.WriteHeader(http.StatusNotFound)
}
