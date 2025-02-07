package internal

import (
	"github.com/kihamo/go-workers"
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/pkg/errors"
)

func (c *Component) LivenessCheck() map[string]dashboard.HealthCheck {
	return map[string]dashboard.HealthCheck{
		"dispatcher": c.DispatcherCheck(),
	}
}

func (c *Component) DispatcherCheck() dashboard.HealthCheck {
	return func() error {
		if c.dispatcher == nil {
			return errors.New("Dispatcher isn't initialized")
		}

		if c.dispatcher.Status() == workers.DispatcherStatusProcess {
			return nil
		}

		return errors.New("Dispatcher status isn't process")
	}
}
