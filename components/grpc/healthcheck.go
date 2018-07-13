package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func HealthCheck(conn *grpc.ClientConn, service string) error {
	healthClient := grpc_health_v1.NewHealthClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	response, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: service,
	})
	if err != nil {
		return err
	}

	if response.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("Server is not healthy. Status is %s", response.Status.String())
	}

	return nil
}
