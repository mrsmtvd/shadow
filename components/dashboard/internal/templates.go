package internal

import (
	"net/url"
	"path"
	"strings"

	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) loadTemplates() error {
	c.renderer = NewRenderer()

	c.renderer.AddGlobalVar("Application", map[string]interface{}{
		"name":       c.application.Name(),
		"version":    c.application.Version(),
		"build":      c.application.Build(),
		"build_date": c.application.BuildDate(),
		"start_date": c.application.StartDate(),
		"uptime":     c.application.Uptime(),
	})
	c.renderer.AddGlobalVar("Config", c.config)

	c.renderer.AddFunc("staticURL", c.funcStaticURL)

	if err := c.renderer.AddBaseLayouts(c.DashboardTemplates()); err != nil {
		return err
	}

	components, err := c.application.GetComponents()
	if err != nil {
		return err
	}

	for _, component := range components {
		if componentTemplate, ok := component.(dashboard.HasTemplates); ok {
			err := c.renderer.AddComponents(component.Name(), componentTemplate.DashboardTemplates())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Component) funcStaticURL(file string, prefix bool) string {
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
