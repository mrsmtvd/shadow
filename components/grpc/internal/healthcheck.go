package internal

import (
	"net"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/grpc"
	"github.com/kihamo/shadow/components/grpc/client"
	g "google.golang.org/grpc"
)

func (c *Component) LivenessCheck() map[string]dashboard.HealthCheck {
	return map[string]dashboard.HealthCheck{
		"server": c.ServerCheck(""),
		"service_grpc": c.ServerCheck("shadow.grpc.Grpc"),
	}
}

func (c *Component) ServerCheck(service string) dashboard.HealthCheck {
	return func() error {
		target := net.JoinHostPort(c.config.String(grpc.ConfigHost), c.config.String(grpc.ConfigPort))
		dial, err := client.DefaultDial(target, g.WithInsecure())
		if err != nil {
			return err
		}

		return grpc.HealthCheck(dial, service)
	}
}
