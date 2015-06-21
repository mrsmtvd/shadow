package frontend

type IndexHandler struct {
	AbstractFrontendHandler
}

func (h *IndexHandler) Handle() {
	h.SetTemplate("index.tpl.html")
	h.View.Context["PageTitle"] = "Home page"
	h.View.Context["Services"] = h.Application.GetServices()
}
