package handlers

import (
	"time"

	"github.com/kihamo/shadow/components/alerts"
	"github.com/kihamo/shadow/components/dashboard"
)

// easyjson:json
type ajaxHandlerResponse struct {
	Icon    string    `json:"icon"`
	Title   string    `json:"title"`
	Message string    `json:"message"`
	Elapsed string    `json:"elapsed"`
	Date    time.Time `json:"date"`
}

type AjaxHandler struct {
	dashboard.Handler

	Component alerts.Component
}

func (h *AjaxHandler) ServeHTTP(w *dashboard.Response, _ *dashboard.Request) {
	list := h.Component.GetAlerts()
	alertsShort := make([]ajaxHandlerResponse, 0, cap(list))

	for i := range list {
		alertsShort = append(alertsShort, ajaxHandlerResponse{
			Icon:    list[i].GetIcon(),
			Title:   list[i].GetTitle(),
			Message: list[i].GetMessage(),
			Elapsed: list[i].GetDateAsMessage(),
			Date:    list[i].GetDate(),
		})
	}

	w.SendJSON(alertsShort)
}
