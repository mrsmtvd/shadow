package frontend

type IndexHandler struct {
	AbstractFrontendHandler
}

func (h *IndexHandler) Handle() {
	h.SetTemplate("index.tpl.html")
	h.View.Context["PageTitle"] = "Application"
	h.View.Context["PageHeader"] = "Application"
	h.View.Context["Services"] = h.Application.GetServices()
	h.View.Context["Resources"] = h.Application.GetResources()
}
