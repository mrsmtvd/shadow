package frontend

type AlertsHandler struct {
	AbstractFrontendHandler
}

func (h *AlertsHandler) Handle() {
	alerts := h.Service.(*FrontendService).GetAlerts()

	if h.IsAjax() {
		alertsShort := make([]map[string]interface{}, 0, cap(alerts))

		for i := range alerts {
			alert := map[string]interface{}{
				"icon":    alerts[i].Icon,
				"title":   alerts[i].Title,
				"message": alerts[i].Message,
				"elapsed": alerts[i].DateAsMessage(),
				"date":    alerts[i].Date,
			}

			alertsShort = append(alertsShort, alert)
		}

		h.SendJSON(alertsShort)
		return
	}

	h.SetTemplate("alerts.tpl.html")
	h.SetPageTitle("Alerts")
	h.SetPageHeader("Alerts")
	h.SetVar("Alerts", alerts)
}
