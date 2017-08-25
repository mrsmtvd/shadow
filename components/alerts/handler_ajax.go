package alerts

import (
	"net/http"
	"time"

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

	component *Component
}

func (h *AjaxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list := h.component.GetAlerts()
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

	dashboard.ResponseFromContext(r.Context()).SendJSON(alertsShort)
}
