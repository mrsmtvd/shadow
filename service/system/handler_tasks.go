package system

import (
	"runtime"

	"github.com/kihamo/shadow/resource"
	"github.com/kihamo/shadow/service/frontend"
)

type TasksHandler struct {
	frontend.AbstractFrontendHandler
}

func (h *TasksHandler) Handle() {
	if h.IsAjax() {
		tasks, _ := h.Application.GetResource("tasks")
		stats := tasks.(*resource.Dispatcher).GetStats()
		stats["goroutines"] = runtime.NumGoroutine()

		h.SendJSON(stats)
		return
	}

	h.SetTemplate("tasks.tpl.html")
	h.View.Context["PageTitle"] = "Tasks"
	h.View.Context["PageHeader"] = "Tasks"
}
