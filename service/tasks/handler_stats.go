package tasks

import (
	"github.com/kihamo/shadow/service/frontend"
)

type StatsHandler struct {
	frontend.AbstractFrontendHandler
}

func (h *StatsHandler) Handle() {
	h.SendJSON(h.Service.(*TasksService).dispatcher.GetStats())
}
