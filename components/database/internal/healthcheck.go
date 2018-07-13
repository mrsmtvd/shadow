package internal

import (
	"context"
	"time"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/database"
	"github.com/pkg/errors"
)

func (c *Component) ReadinessCheck() map[string]dashboard.HealthCheck {
	return map[string]dashboard.HealthCheck{
		"master":     c.MasterExecutorCheck(),
		"slaves":     c.SlavesExecutorCheck(),
		"migrations": c.MigrationsCheck(),
	}
}

func (c *Component) MasterExecutorCheck() dashboard.HealthCheck {
	return func() error {
		s := c.Storage()
		if s == nil {
			return errors.New("Storage isn't initialized")
		}

		return ExecutorCheck(s.Master())
	}
}

func (c *Component) SlavesExecutorCheck() dashboard.HealthCheck {
	return func() error {
		s := c.Storage()
		if s == nil {
			return errors.New("Storage isn't initialized")
		}

		for _, executor := range s.Slaves() {
			if err := ExecutorCheck(executor); err != nil {
				return err
			}
		}

		return nil
	}
}

func (c *Component) MigrationsCheck() dashboard.HealthCheck {
	return func() error {
		c.mutex.RLock()
		defer c.mutex.RUnlock()

		if c.migrationsIsUp {
			return nil
		}

		return c.migrationsError
	}
}

func ExecutorCheck(executor database.Executor) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return executor.Ping(ctx)
}
