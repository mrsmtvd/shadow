package resource

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/GeertJohan/go.rice"
	"github.com/dropbox/godropbox/errors"
	"github.com/kihamo/shadow"
)

type ServiceTemplate interface {
	shadow.Service
	GetTemplateBox() *rice.Box
}

type Template struct {
	Application *shadow.Application
	Globals     map[string]interface{}
}

type TemplateView struct {
	name     string
	service  ServiceTemplate
	template *template.Template
	Globals  map[string]interface{}
	Context  map[string]interface{}
}

func (r *Template) GetName() string {
	return "template"
}

func (r *Template) Init(a *shadow.Application) error {
	r.Application = a
	r.Globals = map[string]interface{}{
		"Application": a,
	}

	return nil
}

func (r *Template) Run() (err error) {
	config, err := r.Application.GetResource("config")
	if err == nil {
		r.Globals["Config"] = config.(*Config)
	}

	return nil
}

func (r *Template) NewView(s string, n string) *TemplateView {
	service, err := r.Application.GetService(s)
	if err != nil {
		panic(fmt.Sprintf("Service \"%s\" not found", s))
	}

	serviceTemplate, ok := service.(ServiceTemplate)
	if !ok {
		panic(fmt.Sprintf("Service \"%s\" not implement ServiceTemplate", s))
	}

	tpl := &TemplateView{
		name:     n,
		service:  serviceTemplate,
		template: template.New(n),
		Globals:  r.Globals,
		Context:  map[string]interface{}{},
	}

	tpl.template.Funcs(template.FuncMap{
		"include": tpl.funcInclude(r),
		"add":     tpl.add,
	})

	return tpl
}

func (v *TemplateView) Render(w io.Writer) error {
	if v.name == "" {
		return nil
	}

	content := v.service.GetTemplateBox().MustString(v.name)
	tpl, err := v.template.Parse(content)
	if err != nil {
		return err
	}

	context := map[string]interface{}{}

	for i := range v.Globals {
		context[i] = v.Globals[i]
	}

	for i := range v.Context {
		context[i] = v.Context[i]
	}

	return tpl.ExecuteTemplate(w, v.name, context)
}

func (v *TemplateView) funcInclude(r *Template) func(args ...string) (string, error) {
	return func(args ...string) (string, error) {
		var (
			name    string
			service string
		)

		if len(args) == 1 {
			name = args[0]
			service = v.service.GetName()
		} else {
			name = args[1]
			service = args[0]
		}

		if name == v.template.Name() {
			return "", errors.New("Recursion include")
		}

		subTemplate := r.NewView(service, name)
		subTemplate.Context = v.Context

		buffer := bytes.NewBuffer([]byte{})

		err := subTemplate.Render(buffer)
		if err != nil {
			return "", err
		}

		return buffer.String(), nil
	}
}

func (v *TemplateView) add(x, y int) (interface{}, error) {
	return x + y, nil
}
