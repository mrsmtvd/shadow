package workers

import (
	"github.com/kihamo/shadow/components/dashboard"
)

type IndexHandler struct {
	dashboard.TemplateHandler
}

func (h *IndexHandler) Handle() {
	h.SetView("workers", "index")
}
