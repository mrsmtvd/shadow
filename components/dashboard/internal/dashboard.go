package internal

import (
	"html/template"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/internal/handlers"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "templates",
	}
}

func (c *Component) DashboardMenu() dashboard.Menu {
	routes := c.DashboardRoutes()

	return dashboard.NewMenu("Dashboard").
		WithRoute(routes[9]).
		WithIcon("dashboard").
		WithChild(dashboard.NewMenu("Components").WithRoute(routes[9])).
		WithChild(dashboard.NewMenu("Environment").WithRoute(routes[3])).
		WithChild(dashboard.NewMenu("Bindata").WithRoute(routes[1])).
		WithChild(dashboard.NewMenu("Routing").WithRoute(routes[4]))
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.NewRoute(
				[]string{http.MethodGet},
				"/"+c.Name()+"/assets/*filepath",
				&assetfs.AssetFS{
					Asset:     Asset,
					AssetDir:  AssetDir,
					AssetInfo: AssetInfo,
					Prefix:    "assets",
				},
				"",
				false),
			dashboard.NewRoute(
				[]string{http.MethodGet},
				"/"+c.Name()+"/bindata",
				&handlers.BindataHandler{
					Application: c.application,
				},
				"",
				true),
			dashboard.NewRoute(
				[]string{http.MethodGet},
				"/"+c.Name()+"/datatables/i18n.json",
				&handlers.DataTablesHandler{
					Application: c.application,
				},
				"",
				false),
			dashboard.NewRoute(
				[]string{http.MethodGet},
				"/"+c.Name()+"/environment",
				&handlers.EnvironmentHandler{},
				"",
				true),
			dashboard.NewRoute(
				[]string{http.MethodGet},
				"/"+c.Name()+"/routing",
				&handlers.RoutingHandler{},
				"",
				true),
			dashboard.NewRoute(
				[]string{http.MethodGet, http.MethodPost},
				dashboard.AuthPath+"/:provider/callback",
				&handlers.AuthHandler{
					IsCallback: true,
				},
				"",
				false),
			dashboard.NewRoute(
				[]string{http.MethodGet, http.MethodPost},
				dashboard.AuthPath+"/:provider",
				&handlers.AuthHandler{},
				"",
				false),
			dashboard.NewRoute(
				[]string{http.MethodGet},
				dashboard.AuthPath,
				&handlers.AuthHandler{},
				"",
				false),
			dashboard.NewRoute(
				[]string{http.MethodGet},
				"/"+c.Name()+"/logout",
				&handlers.LogoutHandler{},
				"",
				true),
		}

		componentsHandler := &handlers.ComponentsHandler{
			Application: c.application,
		}

		c.routes = append(c.routes, []dashboard.Route{
			dashboard.NewRoute(
				[]string{http.MethodGet},
				"/"+c.Name()+"/components",
				componentsHandler,
				"",
				true),
			dashboard.NewRoute(
				[]string{http.MethodGet},
				"/"+c.Name()+"/",
				componentsHandler,
				"",
				true),
		}...)
	}

	return c.routes
}

func (c *Component) DashboardTemplateFunctions() map[string]interface{} {
	return template.FuncMap{
		"raw":        templateFunctionRaw,
		"add":        templateFunctionAdd,
		"mod":        templateFunctionMod,
		"replace":    templateFunctionReplace,
		"staticHTML": templateFunctionStaticHTML,
		"staticURL":  c.templateFunctionStaticURL,
		"date_since": shadow.DateSinceAsMessage,
	}
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
