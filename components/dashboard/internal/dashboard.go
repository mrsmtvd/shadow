package internal

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/Masterminds/sprig/v3"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/internal/handlers"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/misc/time"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	show := func(r *dashboard.Request) bool {
		return r.Config().Bool(config.ConfigDebug)
	}

	return dashboard.NewMenu("Dashboard").
		WithURL("/" + c.Name() + "/components").
		WithIcon("tachometer-alt").
		WithChild(dashboard.NewMenu("Components").WithURL("/" + c.Name() + "/components")).
		WithChild(dashboard.NewMenu("Dependencies").WithURL("/" + c.Name() + "/dependencies")).
		WithChild(dashboard.NewMenu("Environment").WithURL("/" + c.Name() + "/environment")).
		WithChild(dashboard.NewMenu("Asset FS").WithURL("/" + c.Name() + "/assetfs")).
		WithChild(dashboard.NewMenu("Routing").WithURL("/" + c.Name() + "/routing")).
		WithChild(dashboard.NewMenu("Session").WithURL("/" + c.Name() + "/session").WithShow(show)).
		WithChild(dashboard.NewMenu("Health check").
			WithChild(dashboard.NewMenu("Liveness").WithURL("/healthcheck/live?full=1")).
			WithChild(dashboard.NewMenu("Readiness").WithURL("/healthcheck/ready?full=1")))
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	routes := []dashboard.Route{
		dashboard.RouteFromAssetFS(c),
		dashboard.NewRoute("/favicon.ico", dashboard.NewAssetsHandlerByPath(c.AssetFS(), "images/favicon.svg")).
			WithMethods([]string{http.MethodGet}),
		dashboard.NewRoute("/"+c.Name()+"/assetfs", handlers.NewAssetFSHandler(c.registryAssetFS, c.application.BuildDate())).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
		dashboard.NewRoute("/"+c.Name()+"/datatables/i18n.json", &handlers.DataTablesHandler{}).
			WithMethods([]string{http.MethodGet}),
		dashboard.NewRoute("/"+c.Name()+"/environment", &handlers.EnvironmentHandler{}).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
		dashboard.NewRoute("/"+c.Name()+"/dependencies", &handlers.DependenciesHandler{}).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
		dashboard.NewRoute("/"+c.Name()+"/routing", handlers.NewRoutingHandler(c.router)).
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
		dashboard.NewRoute("/"+c.Name()+"/session", &handlers.SessionHandler{}).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
		dashboard.NewRoute("/healthcheck/:healthcheck", handlers.NewHealthCheckHandler(c.components, metricHealthCheckStatus)).
			WithMethods([]string{http.MethodGet}),
	}

	componentsHandler := handlers.NewComponentsHandler(c.application)

	routes = append(routes, []dashboard.Route{
		dashboard.NewRoute("/"+c.Name()+"/components", componentsHandler).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
		dashboard.NewRoute("/"+c.Name()+"/", componentsHandler).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
	}...)

	return routes
}

func (c *Component) DashboardTemplateFunctions() map[string]interface{} {
	list := sprig.FuncMap()
	list["i18n"] = templateFunctionMock(0)
	list["i18nPlural"] = templateFunctionMock(0)
	list["raw"] = templateFunctionRaw
	list["staticHTML"] = templateFunctionStaticHTML
	list["staticURL"] = c.templateFunctionStaticURL
	list["toolbar"] = c.templateFunctionToolbar
	list["date_since"] = time.DateSinceAsMessage
	list["pointer"] = templateFunctionPointer
	list["format_float"] = strconv.FormatFloat

	return list
}

func (c *Component) DashboardToolbar(ctx context.Context) string {
	content, err := c.renderer.RenderLayoutAndReturn(ctx, c.Name(), "toolbar", "blank", nil)

	if err != nil {
		logging.Log(ctx).Error("Failed render toolbar", "error", err.Error())
	}

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
	defer buf.Reset()

	buf.WriteString(c.DashboardToolbar(ctx))

	for _, component := range c.components {
		if component == c {
			continue
		}

		if componentToolbar, ok := component.(dashboard.HasToolbar); ok {
			buf.WriteString(componentToolbar.DashboardToolbar(ctx))
		}
	}

	return template.HTML(buf.String())
}

func templateFunctionMock(i int) func(...interface{}) string {
	return func(args ...interface{}) string {
		if len(args) > i {
			return fmt.Sprintf("%v", args[i])
		}

		return ""
	}
}

func templateFunctionRaw(value interface{}) template.HTML {
	return template.HTML(fmt.Sprint(value))
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

func templateFunctionPointer(v interface{}) interface{} {
	if ref := reflect.ValueOf(v); ref.Kind() == reflect.Ptr {
		if !ref.Elem().IsValid() {
			return nil
		}

		return ref.Elem().Interface()
	}

	return v
}
