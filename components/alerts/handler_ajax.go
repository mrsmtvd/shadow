package alerts

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
)

type AjaxHandler struct {
	dashboard.Handler

	component *Component
}

func (h *AjaxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list := h.component.GetAlerts()
	alertsShort := make([]map[string]interface{}, 0, cap(list))

	for i := range list {
		alert := map[string]interface{}{
			"icon":    list[i].GetIcon(),
			"title":   list[i].GetTitle(),
			"message": list[i].GetMessage(),
			"elapsed": list[i].GetDateAsMessage(),
			"date":    list[i].GetDate(),
		}

		alertsShort = append(alertsShort, alert)
	}

	h.SendJSON(alertsShort, w)
}
