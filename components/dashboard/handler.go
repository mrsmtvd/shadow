package dashboard

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kihamo/gotypes"
)

type Handler interface {
	SetRequest(*http.Request)
	SetResponse(http.ResponseWriter)
	Handle()
}

type HandlerTemplate interface {
	SetRenderer(*Renderer)
	SetView(string, string)
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

	response http.ResponseWriter
	request  *http.Request
}

func (h *AbstractFrontendHandler) GetMethods() []string {
	return []string{http.MethodGet}
}

func (h *AbstractFrontendHandler) IsAuth() bool {
	return true
}

func (h *AbstractFrontendHandler) Request() *http.Request {
	return h.request
}

func (h *AbstractFrontendHandler) SetRequest(r *http.Request) {
	h.request = r
}

func (h *AbstractFrontendHandler) Response() http.ResponseWriter {
	return h.response
}

func (h *AbstractFrontendHandler) SetResponse(r http.ResponseWriter) {
	h.response = r
}

func (h *AbstractFrontendHandler) IsGet() bool {
	return h.request.Method == http.MethodGet
}

func (h *AbstractFrontendHandler) IsPost() bool {
	return h.request.Method == http.MethodPost
}

func (h *AbstractFrontendHandler) IsAjax() bool {
	return h.request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (h *AbstractFrontendHandler) Redirect(location string, code int) {
	http.Redirect(h.response, h.request, location, code)
}

type TemplateHandler struct {
	AbstractFrontendHandler

	renderer  *Renderer
	vars      map[string]interface{}
	component string
	view      string
}

func (h *TemplateHandler) SetResponse(r http.ResponseWriter) {
	h.AbstractFrontendHandler.SetResponse(r)

	r.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (h *TemplateHandler) SetVar(name string, value interface{}) {
	if h.vars == nil {
		h.vars = map[string]interface{}{}
	}

	h.vars[name] = value
}

func (h *TemplateHandler) SetRequest(r *http.Request) {
	h.AbstractFrontendHandler.SetRequest(r)

	h.SetVar("Request", r)
}

func (h *TemplateHandler) SetRenderer(r *Renderer) {
	h.renderer = r
}

func (h *TemplateHandler) SetView(c, v string) {
	h.component = c
	h.view = v
}

func (h *TemplateHandler) Render() {
	if h.renderer != nil {
		err := h.renderer.Render(h.Response(), h.component, h.view, h.vars)
		if err != nil {
			panic(err.Error())
		}
	}
}

type JSONHandler struct {
	AbstractFrontendHandler
}

func (h *JSONHandler) SetResponse(r http.ResponseWriter) {
	h.AbstractFrontendHandler.SetResponse(r)

	r.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (h *JSONHandler) SendJSON(reply interface{}) {
	response, err := json.Marshal(reply)
	if err != nil {
		panic(err.Error())
	}

	h.Response().Write(response)
}

func (h *JSONHandler) DecodeJSON(output interface{}) error {
	decoder := json.NewDecoder(h.Request().Body)

	var in interface{}
	err := decoder.Decode(&in)

	if err != nil {
		return err
	}

	converter := gotypes.NewConverter(in, &output)

	if !converter.Valid() {
		return errors.New("Convert failed")
	}

	return nil
}
