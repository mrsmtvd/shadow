package internal

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/internal/handlers"
	"github.com/kihamo/shadow/misc/time"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	routes := c.DashboardRoutes()

	return dashboard.NewMenu("Dashboard").
		WithRoute(routes[9]).
		WithIcon("dashboard").
		WithChild(dashboard.NewMenu("Components").WithRoute(routes[10])).
		WithChild(dashboard.NewMenu("Environment").WithRoute(routes[3])).
		WithChild(dashboard.NewMenu("Bindata").WithRoute(routes[1])).
		WithChild(dashboard.NewMenu("Routing").WithRoute(routes[4]))
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.RouteFromAssetFS(c),
			dashboard.NewRoute("/"+c.Name()+"/bindata", &handlers.BindataHandler{}).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/datatables/i18n.json", &handlers.DataTablesHandler{}).
				WithMethods([]string{http.MethodGet}),
			dashboard.NewRoute("/"+c.Name()+"/environment", &handlers.EnvironmentHandler{}).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/routing", &handlers.RoutingHandler{}).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute(dashboard.AuthPath+"/:provider/callback", &handlers.AuthHandler{
				IsCallback: true,
			}).
				WithMethods([]string{http.MethodGet, http.MethodPost}),
			dashboard.NewRoute(dashboard.AuthPath+"/:provider", &handlers.AuthHandler{}).
				WithMethods([]string{http.MethodGet, http.MethodPost}),
			dashboard.NewRoute(dashboard.AuthPath, &handlers.AuthHandler{}).
				WithMethods([]string{http.MethodGet}),
			dashboard.NewRoute("/"+c.Name()+"/logout", &handlers.LogoutHandler{}).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/healthcheck/:healthcheck", handlers.NewHealthCheckHandler(c.application, metricHealthCheckStatus)).
				WithMethods([]string{http.MethodGet}),
		}

		componentsHandler := &handlers.ComponentsHandler{}

		c.routes = append(c.routes, []dashboard.Route{
			dashboard.NewRoute("/"+c.Name()+"/components", componentsHandler).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/", componentsHandler).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
		}...)
	}

	return c.routes
}

func (c *Component) DashboardTemplateFunctions() map[string]interface{} {
	return template.FuncMap{
		"i18n":       templateFunctionMock(0),
		"i18nPlural": templateFunctionMock(0),
		"raw":        templateFunctionRaw,
		"add":        templateFunctionAdd,
		"mod":        templateFunctionMod,
		"replace":    templateFunctionReplace,
		"staticHTML": templateFunctionStaticHTML,
		"staticURL":  c.templateFunctionStaticURL,
		"toolbar":    c.templateFunctionToolbar,
		"date_since": time.DateSinceAsMessage,
	}
}

func (c *Component) DashboardToolbar(ctx context.Context) string {
	content, _ := c.renderer.RenderLayoutAndReturn(ctx, c.Name(), "toolbar", "blank", nil)
	return content
}

func (c *Component) templateFunctionStaticURL(file string, prefix bool) string {
	if file == "" {
		return file
	}

	u, err := url.Parse(file)
	if err != nil {
		return file
	}

	if c.application.Build() != "" {
		values := u.Query()
		values.Add("v", c.application.Build())

		u.RawQuery = values.Encode()
	}

	if prefix {
		ext := path.Ext(u.Path)
		lowerExt := strings.ToLower(ext)

		if c.config.Bool(dashboard.ConfigFrontendMinifyEnabled) && (lowerExt == ".css" || lowerExt == ".js") {
			u.Path = u.Path[0:len(u.Path)-len(ext)] + ".min" + ext
		}
	}

	return u.String()
}

func (c *Component) templateFunctionToolbar(opts ...interface{}) template.HTML {
	components, err := c.application.GetComponents()
	if err != nil {
		return ""
	}

	ctx := context.Background()

	if len(opts) > 0 {
		if templateCtx, ok := opts[0].(map[string]interface{}); ok {
			if requestCtx, ok := templateCtx["Request"]; ok {
				if request, ok := requestCtx.(*dashboard.Request); ok {
					ctx = request.Context()
				}
			}
		}
	}

	buf := bytes.NewBuffer(nil)

	for _, component := range components {
		if componentToolbar, ok := component.(dashboard.HasToolbar); ok {
			buf.WriteString(componentToolbar.DashboardToolbar(ctx))
		}
	}

	return template.HTML(buf.String())
}

func templateFunctionStub(...interface{}) string {
	return ""
}

func templateFunctionMock(i int) func(...interface{}) string {
	return func(args ...interface{}) string {
		if len(args) > i {
			return fmt.Sprintf("%v", args[i])
		}

		return ""
	}
}

func templateFunctionRaw(x string) template.HTML {
	return template.HTML(x)
}

func templateFunctionAdd(x, y int) (interface{}, error) {
	return x + y, nil
}

func templateFunctionMod(x, y int) (bool, error) {
	return x%y == 0, nil
}

func templateFunctionReplace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

func templateFunctionStaticHTML(file string) template.HTML {
	if file == "" {
		return template.HTML(file)
	}

	u, err := url.Parse(file)
	if err != nil {
		return template.HTML(file)
	}

	ext := path.Ext(u.Path)

	switch strings.ToLower(ext) {
	case ".css":
		return template.HTML("<link href=\"" + file + "\" rel=\"stylesheet\">")
	case ".js":
		return template.HTML("<script src=\"" + file + "\"></script>")
	}

	return template.HTML(file)
}
