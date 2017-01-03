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

	for _, component := range c.application.GetComponents() {
		if componentTemplate, ok := component.(hasTemplate); ok {
			err := c.renderer.AddComponents(component.GetName(), componentTemplate.GetTemplates())
			if err != nil {
				return err
			}
		}
	}

	c.renderer.AddGlobalVar("Application", c.application)
	c.renderer.AddGlobalVar("Config", c.config)
	c.renderer.AddGlobalVar("AlertsEnabled", c.application.HasComponent("alerts"))

	return nil
}
