package internal

import (
	"github.com/mrsmtvd/shadow/components/dashboard"
)

func (c *Component) initTemplates() error {
	c.renderer.AddGlobalVar("Application", map[string]interface{}{
		"name":       c.application.Name(),
		"version":    c.application.Version(),
		"build":      c.application.Build(),
		"build_date": c.application.BuildDate(),
		"start_date": c.application.StartDate(),
		"uptime":     c.application.Uptime(),
	})

	for name, fn := range c.DashboardTemplateFunctions() {
		c.renderer.AddFunc(name, fn)
	}

	for _, component := range c.components {
		if component == c {
			continue
		}

		if componentTemplateFuncs, ok := component.(dashboard.HasTemplateFunctions); ok {
			for name, fn := range componentTemplateFuncs.DashboardTemplateFunctions() {
				c.renderer.AddFunc(name, fn)
			}
		}
	}

	if err := c.renderer.AddRootTemplates(c.DashboardTemplates()); err != nil {
		return err
	}

	for _, component := range c.components {
		if componentTemplate, ok := component.(dashboard.HasTemplates); ok {
			err := c.renderer.RegisterNamespace(component.Name(), componentTemplate.DashboardTemplates())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
