package dashboard

import (
	"github.com/elazarl/go-bindata-assetfs"
)

type hasTemplate interface {
	GetTemplates() *assetfs.AssetFS
}

func (c *Component) loadTemplates() error {
	c.renderer = NewRenderer()

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

	return nil
}
