package template

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"text/template"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
)

type ServiceTemplate interface {
	shadow.Service
	GetTemplates() *assetfs.AssetFS
}

type Resource struct {
	application *shadow.Application
	Globals     map[string]interface{}
}

type View struct {
	name     string
	service  ServiceTemplate
	template *template.Template
	Globals  map[string]interface{}
	Context  map[string]interface{}
}

func (r *Resource) GetName() string {
	return "template"
}

func (r *Resource) Init(a *shadow.Application) error {
	r.application = a
	r.Globals = map[string]interface{}{
		"Application": a,
	}

	return nil
}

func (r *Resource) Run() (err error) {
	conf, err := r.application.GetResource("config")
	if err == nil {
		r.Globals["Config"] = conf.(*config.Resource).GetAllValues()
	}

	return nil
}

func (r *Resource) NewView(s string, n string) *View {
	service, err := r.application.GetService(s)
	if err != nil {
		panic(fmt.Sprintf("Service \"%s\" not found", s))
	}

	serviceTemplate, ok := service.(ServiceTemplate)
	if !ok {
		panic(fmt.Sprintf("Service \"%s\" not implement ServiceTemplate", s))
	}

	tpl := &View{
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

func (v *View) Render(w io.Writer) error {
	if v.name == "" {
		return nil
	}

	file, err := v.service.GetTemplates().Open(v.name)
	if err != nil {
		return err
	}

	if _, ok := file.(*assetfs.AssetDirectory); ok {
		return errors.New(v.name + " is directory")
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(file.(*assetfs.AssetFile).Reader)

	tpl, err := v.template.Parse(buf.String())
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

func (v *View) funcInclude(r *Resource) func(args ...string) (string, error) {
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

func (v *View) add(x, y int) (interface{}, error) {
	return x + y, nil
}
