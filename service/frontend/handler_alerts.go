package frontend

import (
	"github.com/kihamo/shadow/resource/alerts"
)

type AlertsHandler struct {
	AbstractFrontendHandler
}

func (h *AlertsHandler) Handle() {
	resourceAlerts, _ := h.Application.GetResource("alerts")
	list := resourceAlerts.(*alerts.Alerts).GetAlerts()

	if h.IsAjax() {
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

		h.SendJSON(alertsShort)
		return
	}

	h.SetTemplate("alerts.tpl.html")
	h.SetPageTitle("Alerts")
	h.SetPageHeader("Alerts")
	h.SetVar("Alerts", list)
}
