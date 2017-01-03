package alerts

import (
	"github.com/kihamo/shadow/components/dashboard"
)

type IndexHandler struct {
	dashboard.TemplateHandler

	component *Component
}

func (h *IndexHandler) Handle() {
	h.SetView("alerts", "index")

	list := h.component.GetAlerts()
	h.SetVar("alerts", list)
}
