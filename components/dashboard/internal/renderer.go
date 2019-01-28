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

type templatesNamespace struct {
	fs        *assetfs.AssetFS
	templates map[string]*template.Template
}

type Renderer struct {
	rootTemplate *template.Template
	globals      map[string]interface{}
	namespaces   map[string]*templatesNamespace
}

func NewRenderer() *Renderer {
	root := template.New(TemplateRootName).
		Funcs(dashboard.DefaultTemplateFunctions.FuncMap())

	r := &Renderer{
		rootTemplate: root,
		globals:      map[string]interface{}{},
		namespaces:   map[string]*templatesNamespace{},
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

func (r *Renderer) IsRegisterNamespace(ns string) bool {
	_, ok := r.namespaces[ns]
	return ok
}

func (r *Renderer) RegisterNamespace(ns string, fs *assetfs.AssetFS) error {
	if r.IsRegisterNamespace(ns) {
		return errors.New("namesapce " + ns + " already exists")
	}

	r.namespaces[ns] = &templatesNamespace{
		fs:        fs,
		templates: make(map[string]*template.Template, 0),
	}

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
	tpl, err := r.getLazyViewTemplate(ns, view)
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

func (r *Renderer) getLazyViewTemplate(ns, view string) (*template.Template, error) {
	namespace, ok := r.namespaces[ns]
	if !ok {
		return nil, errors.New("namespace \"" + ns + "\" not found")
	}

	view += TemplatePostfix

	tpl, ok := namespace.templates[view]
	if ok {
		return tpl, nil
	}

	files, err := r.getTemplateFiles(TemplateViewsDir, namespace.fs)
	if err != nil {
		return nil, err
	}

	var (
		found   bool
		content []byte
	)

	for name, body := range files {
		if name == view {
			found = true
			content = body
			break
		}
	}

	if !found {
		return nil, errors.New("template \"" + view + "\" for namespace \"" + ns + "\" not found")
	}

	// layouts
	tpl, err = r.rootTemplate.Clone()
	if err != nil {
		return nil, err
	}

	if files, err := r.getTemplateFiles(TemplateLayoutsDir, namespace.fs); err == nil {
		for layout, body := range files {
			layout = strings.TrimSuffix(layout, TemplatePostfix)

			t := tpl.Lookup(layout)
			if t == nil {
				t = tpl.New(layout)
			}

			if _, err := t.Parse(string(body)); err != nil {
				return nil, err
			}
		}
	}

	// views
	if tpl, err = tpl.Parse(string(content)); err != nil {
		return nil, err
	}

	namespace.templates[view] = tpl
	return tpl, nil
}
