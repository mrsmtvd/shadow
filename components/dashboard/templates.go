package dashboard

import (
	"net/url"
	"path"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/config"
)

type hasTemplate interface {
	GetTemplates() *assetfs.AssetFS
}

func (c *Component) loadTemplates() error {
	c.renderer = NewRenderer()

	c.renderer.AddGlobalVar("Application", map[string]interface{}{
		"name":       c.application.GetName(),
		"version":    c.application.GetVersion(),
		"build":      c.application.GetBuild(),
		"build_date": c.application.GetBuildDate(),
		"start_date": c.application.GetStartDate(),
		"uptime":     c.application.GetUptime(),
	})
	c.renderer.AddGlobalVar("Config", c.config.GetAllValues())
	c.renderer.AddGlobalVar("AlertsEnabled", c.application.HasComponent("alerts"))

	c.renderer.AddFunc("staticURL", c.funcStaticURL)

	if err := c.renderer.AddBaseLayouts(c.GetTemplates()); err != nil {
		return err
	}

	components, err := c.application.GetComponents()
	if err != nil {
		return err
	}

	for _, component := range components {
		if componentTemplate, ok := component.(hasTemplate); ok {
			err := c.renderer.AddComponents(component.GetName(), componentTemplate.GetTemplates())
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

	if c.application.GetBuild() != "" {
		values := u.Query()
		values.Add("v", c.application.GetBuild())

		u.RawQuery = values.Encode()
	}

	if prefix {
		ext := path.Ext(u.Path)
		lowerExt := strings.ToLower(ext)

		if !c.config.GetBool(config.ConfigDebug) && (lowerExt == ".css" || lowerExt == ".js") {
			u.Path = u.Path[0:len(u.Path)-len(ext)] + ".min" + ext
		}
	}

	return u.String()
}
