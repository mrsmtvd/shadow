package system

import (
	"github.com/kihamo/shadow/resource"
	"github.com/kihamo/shadow/service/frontend"
)

type TasksHandler struct {
	frontend.AbstractFrontendHandler
}

func (h *TasksHandler) Handle() {
	if h.IsAjax() {
		tasks, _ := h.Application.GetResource("tasks")
		h.SendJSON(tasks.(*resource.Dispatcher).GetStats())
		return
	}

	h.SetTemplate("tasks.tpl.html")
	h.View.Context["PageTitle"] = "Tasks"
}
