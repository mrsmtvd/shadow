package internal

import (
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

const (
	TemplatePostfix       = ".html"
	TemplateLayoutsDir    = "templates/layouts"
	TemplateViewsDir      = "templates/views"
	TemplateDefaultLayout = "base"
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
		funcs:       dashboard.DefaultTemplateFunctions.FuncMap(),
	}

	return r
}

func (r *Renderer) AddFunc(name string, f interface{}) {
	r.funcs[name] = f
}

func (r *Renderer) AddBaseLayouts(fs *assetfs.AssetFS) error {
	files, err := r.getTemplateFiles(TemplateLayoutsDir, fs)
	if err != nil {
		return err
	}

	for name, content := range files {
		r.baseLayouts[strings.TrimSuffix(name, TemplatePostfix)] = content
	}

	return nil
}

func (r *Renderer) AddGlobalVar(key string, value interface{}) {
	r.globals[key] = value
}

func (r *Renderer) AddComponents(componentName string, fs *assetfs.AssetFS) error {
	baseComponent := template.New("_component").Funcs(r.funcs)

	// layouts
	for name, content := range r.baseLayouts {
		baseComponent.New(name).Parse(content)
	}

	if files, err := r.getTemplateFiles(TemplateLayoutsDir, fs); err == nil {
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
	files, err := r.getTemplateFiles(TemplateViewsDir, fs)
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

	r.views[componentName] = views

	return nil
}

func (r *Renderer) Render(ctx context.Context, componentName, viewName string, data map[string]interface{}) error {
	return r.RenderLayout(ctx, componentName, viewName, TemplateDefaultLayout, data)
}

func (r *Renderer) RenderLayout(ctx context.Context, componentName, viewName, layoutName string, data map[string]interface{}) error {
	component, ok := r.views[componentName]
	if !ok {
		return fmt.Errorf("DashboardTemplates for component \"%s\" not found", componentName)
	}

	view, ok := component[viewName+TemplatePostfix]
	if !ok {
		return fmt.Errorf("Template \"%s\" for component \"%s\" not found", viewName, componentName)
	}

	executeData := r.getContextVariables(ctx)
	executeData["ComponentName"] = componentName
	executeData["ViewName"] = viewName
	executeData["LayoutName"] = layoutName

	for i := range r.globals {
		executeData[i] = r.globals[i]
	}

	if data != nil {
		for i := range data {
			executeData[i] = data[i]
		}
	}

	return view.ExecuteTemplate(dashboard.ResponseFromContext(ctx), layoutName, executeData)
}

func (r *Renderer) getContextVariables(ctx context.Context) map[string]interface{} {
	request := dashboard.RequestFromContext(ctx)

	return map[string]interface{}{
		"Request": request,
		"User":    request.User(),
	}
}

func (r *Renderer) getTemplateFiles(directory string, f *assetfs.AssetFS) (map[string]string, error) {
	files, err := f.AssetDir(directory)
	if err != nil {
		return nil, err
	}

	templates := make(map[string]string, 0)

	for _, file := range files {
		if !strings.HasSuffix(file, TemplatePostfix) {
			continue
		}

		content, err := f.Asset(directory + "/" + file)
		if err != nil {
			continue
		}

		templates[file] = string(content)
	}

	return templates, nil
}
