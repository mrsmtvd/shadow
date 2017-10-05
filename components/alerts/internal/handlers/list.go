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
	if r.IsAjax() {
		list := h.Component.GetAlerts()
		alertsShort := make([]listHandlerResponse, 0, cap(list))

		for i := range list {
			alertsShort = append(alertsShort, listHandlerResponse{
				Icon:    list[i].GetIcon(),
				Title:   list[i].GetTitle(),
				Message: list[i].GetMessage(),
				Elapsed: list[i].GetDateAsMessage(),
				Date:    list[i].GetDate(),
			})
		}

		w.SendJSON(alertsShort)
		return
	}

	h.Render(r.Context(), h.Component.GetName(), "list", map[string]interface{}{
		"alerts": h.Component.GetAlerts(),
	})
}
