package internal

import (
	"errors"

	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) ReadinessCheck() map[string]dashboard.HealthCheck {
	// TODO: check shutdown application

	hc := make(map[string]dashboard.HealthCheck, len(c.components))

	for _, cmp := range c.components {
		hc["component_"+cmp.Name()+"_ready"] = c.ComponentReadyCheck(cmp.Name())
	}

	return hc
}

func (c *Component) ComponentReadyCheck(name string) dashboard.HealthCheck {
	return func() error {
		if !c.application.IsReadyComponent(name) {
			return errors.New("not ready")
		}

		return nil
	}
}
