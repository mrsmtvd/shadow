package handlers

import (
	"time"

	"github.com/kihamo/shadow/components/alerts"
	"github.com/kihamo/shadow/components/dashboard"
)

// easyjson:json
type listHandlerResponse struct {
	Icon    string    `json:"icon"`
	Title   string    `json:"title"`
	Message string    `json:"message"`
	Elapsed string    `json:"elapsed"`
	Date    time.Time `json:"date"`
}

type ListHandler struct {
	dashboard.Handler

	Component alerts.Component
}

func (h *ListHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	list := h.Component.GetAlerts()

	if r.IsAjax() {
		reply := make([]listHandlerResponse, 0, len(list))

		for i := range list {
			reply = append(reply, listHandlerResponse{
				Icon:    list[i].Icon(),
				Title:   list[i].Title(),
				Message: list[i].Message(),
				Elapsed: list[i].DateAsMessage(),
				Date:    list[i].Date(),
			})
		}

		w.SendJSON(reply)
		return
	}

	h.Render(r.Context(), h.Component.GetName(), "list", map[string]interface{}{
		"alerts": list,
	})
}
