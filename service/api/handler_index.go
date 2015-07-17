package api

import (
	"fmt"

	"github.com/kihamo/shadow/service/frontend"
)

type IndexHandler struct {
	frontend.AbstractFrontendHandler
}

func (h *IndexHandler) Handle() {
	h.SetTemplate("index.tpl.html")
	h.View.Context["PageTitle"] = "Api"

	service := h.Service.(*ApiService)

	host := service.config.GetString("api-host")
	if host == "0.0.0.0" {
		host = "localhost"
	}

	h.View.Context["ApiUrl"] = fmt.Sprintf("ws://%s:%d/", host, service.config.GetInt64("api-port"))
}
