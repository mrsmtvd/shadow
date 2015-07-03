package frontend

import (
	"encoding/json"
	"net/http"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
)

type Handler interface {
	http.Handler

	Init(*shadow.Application, shadow.Service)
	InitRequest(http.ResponseWriter, *http.Request)
	SetTemplate(string)
	Handle()
	Render()
}

type AbstractFrontendHandler struct {
	Handler

	Application *shadow.Application
	Service     shadow.Service
	Template    *resource.Template
	View        *resource.TemplateView

	Output http.ResponseWriter
	Input  *http.Request
}

func (h *AbstractFrontendHandler) Init(a *shadow.Application, s shadow.Service) {
	templateResource, _ := a.GetResource("template")

	h.Application = a
	h.Service = s
	h.Template = templateResource.(*resource.Template)
}

func (h *AbstractFrontendHandler) InitRequest(out http.ResponseWriter, in *http.Request) {
	h.Output = out
	h.Input = in

	h.View = nil

	out.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (h *AbstractFrontendHandler) SetTemplate(name string) {
	h.View = h.Template.NewView(h.Service.GetName(), name)
	h.View.Context["Request"] = h.Input
}

func (h *AbstractFrontendHandler) Render() {
	if h.View != nil {
		err := h.View.Render(h.Output)
		if err != nil {
			panic(err.Error())
		}
	}
}

func (h *AbstractFrontendHandler) IsAjax() bool {
	return h.Input.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (h *AbstractFrontendHandler) SendJSON(reply interface{}) {
	response, err := json.Marshal(reply)
	if err != nil {
		panic(err.Error())
	}

	h.Output.Header().Set("Content-Type", "application/json; charset=utf-8")
	h.Output.Write(response)
}
