package frontend

type IndexHandler struct {
	AbstractFrontendHandler
}

func (h *IndexHandler) Handle() {
	h.SetTemplate("index.tpl.html")
	h.SetPageTitle("Application")
	h.SetPageHeader("Application")
	h.SetVar("Services", h.Application.GetServices())
	h.SetVar("Resources", h.Application.GetResources())
}
