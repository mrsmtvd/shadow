package internal

import (
	"net"
	"os"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/grpc"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/metrics"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logger.Logger
	server      *g.Server
	routes      []dashboard.Route
}

func (c *Component) GetName() string {
	return grpc.ComponentName
}

func (c *Component) GetVersion() string {
	return grpc.ComponentVersion + "/" + g.Version
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
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) error {
	c.logger = logger.NewOrNop(c.GetName(), c.application)

	var serverOptions []g.ServerOption

	// interceptors
	unaryInterceptors := []g.UnaryServerInterceptor{
		NewConfigUnaryServerInterceptor(c.config),
		NewLoggerUnaryServerInterceptor(c.logger),
	}
	streamInterceptors := []g.StreamServerInterceptor{
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

	c.server = g.NewServer(serverOptions...)

	components, err := c.application.GetComponents()
	if err != nil {
		return err
	}

	for _, cmp := range components {
		if cmpGrpc, ok := cmp.(grpc.HasGrpcServer); ok {
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

		if c.config.Bool(grpc.ConfigReflectionEnabled) {
			reflection.Register(c.server)
		}

		addr := net.JoinHostPort(c.config.String(grpc.ConfigHost), c.config.String(grpc.ConfigPort))
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
