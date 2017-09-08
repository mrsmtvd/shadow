package grpc

import (
	"net"
	"os"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	ComponentName = "grpc"
)

type GrpcService interface {
	RegisterGrpcServer(*grpc.Server)
}

type Component struct {
	application shadow.Application
	config      *config.Component
	logger      logger.Logger
	server      *grpc.Server
}

func (c *Component) GetName() string {
	return ComponentName
}

func (c *Component) GetVersion() string {
	return ComponentVersion
}

func (c *Component) GetDependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: logger.ComponentName,
		},
		{
			Name: metrics.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(*config.Component)

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) error {
	c.logger = logger.NewOrNop(c.GetName(), c.application)

	var serverOptions []grpc.ServerOption

	// interceptors
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		NewConfigUnaryServerInterceptor(c.config),
		NewLoggerUnaryServerInterceptor(c.logger),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		NewConfigStreamServerInterceptor(c.config),
		NewLoggerStreamServerInterceptor(c.logger),
	}

	if c.application.HasComponent(metrics.ComponentName) {
		unaryInterceptors = append(unaryInterceptors, NewMetricsUnaryServerInterceptor())
		streamInterceptors = append(streamInterceptors, NewMetricsStreamServerInterceptor())
	}

	unaryInterceptors = append(unaryInterceptors, grpc_recovery.UnaryServerInterceptor())
	streamInterceptors = append(streamInterceptors, grpc_recovery.StreamServerInterceptor())

	serverOptions = append(serverOptions, grpc_middleware.WithUnaryServerChain(unaryInterceptors...))
	serverOptions = append(serverOptions, grpc_middleware.WithStreamServerChain(streamInterceptors...))

	c.server = grpc.NewServer(serverOptions...)

	components, err := c.application.GetComponents()
	if err != nil {
		return err
	}

	for _, cmp := range components {
		if cmpGrpc, ok := cmp.(GrpcService); ok {
			cmpGrpc.RegisterGrpcServer(c.server)
		}
	}

	for service, info := range c.server.GetServiceInfo() {
		for _, method := range info.Methods {
			c.logger.Debugf("Add method /%s/%s", service, method.Name)
		}
	}

	var wgListen sync.WaitGroup
	wgListen.Add(1)

	go func() {
		defer wg.Done()

		if c.config.GetBool(ConfigReflectionEnabled) {
			reflection.Register(c.server)
		}

		addr := net.JoinHostPort(c.config.GetString(ConfigHost), c.config.GetString(ConfigPort))
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			c.logger.Fatalf("Failed to listen [%d]: %s\n", os.Getpid(), err.Error())
		}

		c.logger.Info("Running service", map[string]interface{}{
			"addr": addr,
			"pid":  os.Getpid(),
		})

		wgListen.Done()

		if err := c.server.Serve(lis); err != nil {
			c.logger.Fatalf("Failed to serve [%d]: %s\n", os.Getpid(), err.Error())
		}
	}()

	wgListen.Wait()
	return nil
}
