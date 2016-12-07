package frontend

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kihamo/gotypes"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/template"
)

const (
	PageTitleVar = "PageTitle"
	PageHeadeVar = "PageHeader"
)

type Handler interface {
	http.Handler

	Init(*shadow.Application, shadow.Service)
	InitRequest(http.ResponseWriter, *http.Request)
	SetTemplate(string)
	Handle()
	Render()
}

type HandlerAuth interface {
	IsAuth() bool
}

type HandlerPanic interface {
	SetError(interface{})
}

type AbstractFrontendHandler struct {
	Handler

	Application *shadow.Application
	Service     shadow.Service
	Template    *template.Resource
	View        *template.View

	Output http.ResponseWriter
	Input  *http.Request
}

func (h *AbstractFrontendHandler) IsAuth() bool {
	return true
}

func (h *AbstractFrontendHandler) Init(a *shadow.Application, s shadow.Service) {
	templateResource, _ := a.GetResource("template")

	h.Application = a
	h.Service = s
	h.Template = templateResource.(*template.Resource)
}

func (h *AbstractFrontendHandler) InitRequest(out http.ResponseWriter, in *http.Request) {
	h.Output = out
	h.Input = in

	h.View = nil

	out.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (h *AbstractFrontendHandler) SetTemplate(name string) {
	h.View = h.Template.NewView(h.Service.GetName(), name)
	h.SetPageHeader("")
	h.SetPageTitle("")
	h.SetVar("Request", h.Input)
}

func (h *AbstractFrontendHandler) SetVar(name string, value interface{}) {
	h.View.Context[name] = value
}

func (h *AbstractFrontendHandler) SetPageTitle(title string) {
	h.SetVar(PageTitleVar, title)
}

func (h *AbstractFrontendHandler) SetPageHeader(header string) {
	h.SetVar(PageHeadeVar, header)
}

func (h *AbstractFrontendHandler) Render() {
	if h.View != nil {
		err := h.View.Render(h.Output)
		if err != nil {
			panic(err.Error())
		}
	}
}

func (h *AbstractFrontendHandler) IsGet() bool {
	return h.Input.Method == "GET"
}

func (h *AbstractFrontendHandler) IsPost() bool {
	return h.Input.Method == "POST"
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

func (h *AbstractFrontendHandler) DecodeJSON(output interface{}) error {
	decoder := json.NewDecoder(h.Input.Body)

	var in interface{}
	err := decoder.Decode(&in)

	if err != nil {
		return err
	}

	converter := gotypes.NewConverter(in, &output)

	if !converter.Valid() {
		return errors.New("Convert fail")
	}

	return nil
}
