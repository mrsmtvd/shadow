package internal

import (
	"net"
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/grpc"
	"github.com/mrsmtvd/shadow/components/grpc/client"
	g "google.golang.org/grpc"
)

const (
	requestTimeout = time.Second * 2
)

func (c *Component) LivenessCheck() map[string]dashboard.HealthCheck {
	return map[string]dashboard.HealthCheck{
		"server":       c.ServerCheck(""),
		"service_grpc": c.ServerCheck("shadow.grpc.Grpc"),
	}
}

func (c *Component) ServerCheck(service string) dashboard.HealthCheck {
	return healthcheck.Timeout(func() error {
		target := net.JoinHostPort(c.config.String(grpc.ConfigHost), c.config.String(grpc.ConfigPort))
		dial, err := client.DefaultDial(target, g.WithInsecure())
		if err != nil {
			return err
		}

		return grpc.HealthCheck(dial, service)
	}, requestTimeout)
}
