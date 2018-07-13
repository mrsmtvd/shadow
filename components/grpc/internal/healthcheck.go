package internal

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/grpc"
	"github.com/kihamo/shadow/components/grpc/client"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func (c *Component) LivenessCheck() map[string]dashboard.HealthCheck {
	return map[string]dashboard.HealthCheck{
		"server": c.ServerCheck(),
	}
}

func (c *Component) ServerCheck() dashboard.HealthCheck {
	return func() error {
		target := net.JoinHostPort(c.config.String(grpc.ConfigHost), c.config.String(grpc.ConfigPort))
		dial, err := client.DefaultDial(target, g.WithInsecure())
		if err != nil {
			return err
		}

		healthClient := grpc_health_v1.NewHealthClient(dial)
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		response, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
		if err != nil {
			return err
		}

		if response.Status != grpc_health_v1.HealthCheckResponse_SERVING {
			return fmt.Errorf("Server is not healthy.Status is %s", response.Status.String())
		}

		return nil
	}
}
