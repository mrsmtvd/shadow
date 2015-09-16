package frontend

type AlertsHandler struct {
	AbstractFrontendHandler
}

func (h *AlertsHandler) Handle() {
	h.SetTemplate("alerts.tpl.html")
	h.View.Context["PageTitle"] = "Alerts"
	h.View.Context["PageHeader"] = "Alerts"
	h.View.Context["Alerts"] = h.Service.(*FrontendService).GetAlerts()
}
