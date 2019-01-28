package internal

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"io"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

const (
	TemplateRootName      = "_root"
	TemplatePostfix       = ".html"
	TemplateLayoutsDir    = "templates/layouts"
	TemplateViewsDir      = "templates/views"
	TemplateDefaultLayout = "base"
)

type Renderer struct {
	rootTemplate *template.Template
	globals      map[string]interface{}
	templates    map[string]map[string]*template.Template
}

func NewRenderer() *Renderer {
	root := template.New(TemplateRootName).
		Funcs(dashboard.DefaultTemplateFunctions.FuncMap())

	r := &Renderer{
		rootTemplate: root,
		globals:      map[string]interface{}{},
		templates:    map[string]map[string]*template.Template{},
	}

	return r
}

func (r *Renderer) AddFunc(name string, f interface{}) {
	r.rootTemplate = r.rootTemplate.Funcs(template.FuncMap{
		name: f,
	})
}

func (r *Renderer) AddRootTemplates(fs *assetfs.AssetFS) error {
	files, err := r.getTemplateFiles(TemplateLayoutsDir, fs)
	if err != nil {
		return err
	}

	tpl := r.rootTemplate

	for layout, content := range files {
		layout = strings.TrimSuffix(layout, TemplatePostfix)

		if tpl, err = tpl.New(layout).Parse(string(content)); err != nil {
			return err
		}
	}

	r.rootTemplate = tpl
	return nil
}

func (r *Renderer) AddGlobalVar(key string, value interface{}) {
	r.globals[key] = value
}

func (r *Renderer) IsRegisterNamespace(componentName string) bool {
	_, ok := r.templates[componentName]
	return ok
}

func (r *Renderer) RegisterNamespace(ns string, fs *assetfs.AssetFS) error {
	layouts, err := r.rootTemplate.Clone()
	if err != nil {
		return err
	}

	// layouts
	if files, err := r.getTemplateFiles(TemplateLayoutsDir, fs); err == nil {
		for layout, content := range files {
			layout = strings.TrimSuffix(layout, TemplatePostfix)

			tpl := layouts.Lookup(layout)
			if tpl == nil {
				tpl = layouts.New(layout)
			}

			if _, err := tpl.Parse(string(content)); err != nil {
				return err
			}
		}
	}

	// views
	files, err := r.getTemplateFiles(TemplateViewsDir, fs)
	if err != nil {
		return nil
	}

	templates := map[string]*template.Template{}
	for view, content := range files {
		tpl, err := layouts.Clone()
		if err != nil {
			return err
		}

		if tpl, err = tpl.Parse(string(content)); err != nil {
			return err
		}

		templates[view] = tpl
	}

	r.templates[ns] = templates

	return nil
}

func (r *Renderer) RenderAndReturn(ctx context.Context, ns, view string, data map[string]interface{}) (string, error) {
	wr := bytes.NewBuffer(nil)
	err := r.Render(wr, ctx, ns, view, data)

	return wr.String(), err
}

func (r *Renderer) Render(wr io.Writer, ctx context.Context, ns, view string, data map[string]interface{}) error {
	return r.RenderLayout(wr, ctx, ns, view, TemplateDefaultLayout, data)
}

func (r *Renderer) RenderLayoutAndReturn(ctx context.Context, ns, view, layout string, data map[string]interface{}) (string, error) {
	wr := bytes.NewBuffer(nil)
	err := r.RenderLayout(wr, ctx, ns, view, layout, data)

	return wr.String(), err
}

func (r *Renderer) RenderLayout(wr io.Writer, ctx context.Context, ns, view, layout string, data map[string]interface{}) error {
	tpl, err := r.getViewTemplate(ns, view)
	if err != nil {
		return err
	}

	executeData := r.getContextVariables(ctx)
	executeData["ComponentName"] = ns
	executeData["ViewName"] = view
	executeData["LayoutName"] = layout

	for i := range r.globals {
		executeData[i] = r.globals[i]
	}

	for i := range data {
		executeData[i] = data[i]
	}

	return tpl.ExecuteTemplate(wr, layout, executeData)
}

func (r *Renderer) getContextVariables(ctx context.Context) map[string]interface{} {
	vars := map[string]interface{}{}

	request := dashboard.RequestFromContext(ctx)
	if request != nil {
		vars["Request"] = request
		vars["User"] = request.User()
	}

	return vars
}

func (r *Renderer) getTemplateFiles(directory string, f *assetfs.AssetFS) (map[string][]byte, error) {
	files, err := f.AssetDir(directory)
	if err != nil {
		return nil, err
	}

	templates := make(map[string][]byte)

	for _, file := range files {
		if !strings.HasSuffix(file, TemplatePostfix) {
			continue
		}

		content, err := f.Asset(directory + "/" + file)
		if err != nil {
			continue
		}

		templates[file] = content
	}

	return templates, nil
}

func (r *Renderer) getViewTemplate(ns, view string) (*template.Template, error) {
	tpls, ok := r.templates[ns]
	if !ok {
		return nil, errors.New("templates for namespace \"" + ns + "\" not found")
	}

	tpl, ok := tpls[view+TemplatePostfix]
	if !ok {
		return nil, errors.New("template \"" + view + "\" for namespace \"" + ns + "\" not found")
	}

	return tpl, nil
}
