package dashboard

import (
	"github.com/kihamo/shadow"
)

type ComponentsHandler struct {
	TemplateHandler

	application *shadow.Application
}

func (h *ComponentsHandler) Handle() {
	h.SetView("dashboard", "components")
	h.SetVar("components", h.application.GetComponents())
}
