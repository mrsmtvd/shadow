package internal

import (
	"errors"

	"github.com/mrsmtvd/shadow"
	"github.com/mrsmtvd/shadow/components/dashboard"
)

func (c *Component) ReadinessCheck() map[string]dashboard.HealthCheck {
	hc := make(map[string]dashboard.HealthCheck, len(c.components))

	for _, cmp := range c.components {
		hc["component_"+cmp.Name()+"_ready"] = c.ComponentReadyCheck(cmp.Name())
	}

	return hc
}

func (c *Component) ComponentReadyCheck(name string) dashboard.HealthCheck {
	return func() error {
		switch c.application.StatusComponent(name) {
		case shadow.ComponentStatusReady, shadow.ComponentStatusFinished:
			return nil
		}

		return errors.New("not ready")
	}
}
