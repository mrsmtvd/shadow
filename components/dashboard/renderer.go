package dashboard

import (
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
)

const (
	TemplatePostfix    = ".html"
	TemplateLayoutsDir = "templates/layouts"
	TemplateViewsDir   = "templates/views"
)

type Renderer struct {
	baseLayouts map[string]string
	globals     map[string]interface{}
	views       map[string]map[string]*template.Template
	funcs       template.FuncMap
}

func NewRenderer() *Renderer {
	r := &Renderer{
		baseLayouts: map[string]string{},
		globals:     map[string]interface{}{},
		views:       map[string]map[string]*template.Template{},
	}

	r.funcs = template.FuncMap{
		"add": r.funcAdd,
		"mod": r.funcMod,
	}

	return r
}

func (r *Renderer) AddBaseLayouts(f *assetfs.AssetFS) error {
	files, err := r.getTemplateFiles(TemplateLayoutsDir, f)
	if err != nil {
		return err
	}

	for name, content := range files {
		r.baseLayouts[strings.TrimSuffix(name, TemplatePostfix)] = content
	}

	return nil
}

func (r *Renderer) AddGlobalVar(n string, v interface{}) {
	r.globals[n] = v
}

func (r *Renderer) AddComponents(c string, f *assetfs.AssetFS) error {
	baseComponent := template.New("_component").Funcs(r.funcs)

	// layouts
	for name, content := range r.baseLayouts {
		baseComponent.New(name).Parse(content)
	}

	if files, err := r.getTemplateFiles(TemplateLayoutsDir, f); err == nil {
		for name, content := range files {
			tplName := strings.TrimSuffix(name, TemplatePostfix)

			tpl := baseComponent.Lookup(tplName)
			if tpl == nil {
				tpl.New(tplName)
			}

			tpl.Parse(content)
		}
	}

	// views
	files, err := r.getTemplateFiles(TemplateViewsDir, f)
	if err != nil {
		return nil
	}

	views := map[string]*template.Template{}
	for name, content := range files {
		view, err := baseComponent.Clone()
		if err != nil {
			return err
		}

		if view, err = view.Parse(content); err != nil {
			return err
		}

		views[name] = view
	}

	r.views[c] = views

	return nil
}

func (r *Renderer) Render(ctx context.Context, c, v string, d map[string]interface{}) error {
	component, ok := r.views[c]
	if !ok {
		return fmt.Errorf("Templates for component \"%s\" not found", c)
	}

	view, ok := component[v+TemplatePostfix]
	if !ok {
		return fmt.Errorf("Template \"%s\" for component \"%s\" not found", v, c)
	}

	context := map[string]interface{}{
		"Request": RequestFromContext(ctx),
	}

	for i := range r.globals {
		context[i] = r.globals[i]
	}

	if d != nil {
		for i := range d {
			context[i] = d[i]
		}
	}

	return view.ExecuteTemplate(ResponseFromContext(ctx), "main", context)
}

func (r *Renderer) getTemplateFiles(d string, f *assetfs.AssetFS) (map[string]string, error) {
	files, err := f.AssetDir(d)
	if err != nil {
		return nil, err
	}

	templates := make(map[string]string, 0)

	for _, file := range files {
		if !strings.HasSuffix(file, TemplatePostfix) {
			continue
		}

		content, err := f.Asset(d + "/" + file)
		if err != nil {
			continue
		}

		templates[file] = string(content)
	}

	return templates, nil
}

func (r *Renderer) funcAdd(x, y int) (interface{}, error) {
	return x + y, nil
}

func (r *Renderer) funcMod(x, y int) (bool, error) {
	return x%y == 0, nil
}
