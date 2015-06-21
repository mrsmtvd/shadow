package aws

import (
	"github.com/kihamo/shadow/service/frontend"
)

type IndexHandler struct {
	frontend.AbstractFrontendHandler
}

func (h *IndexHandler) Handle() {
	h.SetTemplate("index.tpl.html")
	h.View.Context["PageTitle"] = "Aws"

	service := h.Service.(*AwsService)
	h.View.Context["Applications"] = service.applications
	h.View.Context["Subscriptions"] = service.subscriptions
	h.View.Context["Topics"] = service.topics
}
