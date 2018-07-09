package internal

import (
	"net"
	"os"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/grpc"
	"github.com/kihamo/shadow/components/grpc/interceptor"
	"github.com/kihamo/shadow/components/grpc/stats"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/metrics"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	s "google.golang.org/grpc/stats"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logger.Logger
	server      *g.Server
	routes      []dashboard.Route
}

func (c *Component) Name() string {
	return grpc.ComponentName
}

func (c *Component) Version() string {
	return grpc.ComponentVersion + "/" + g.Version
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: i18n.ComponentName,
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
	c.logger = logger.NewOrNop(c.Name(), c.application)
	grpclog.SetLoggerV2(grpc.NewLogger(c.logger))

	var serverOptions []g.ServerOption

	// interceptors
	unaryInterceptors := []g.UnaryServerInterceptor{
		interceptor.NewRecoverUnaryServerInterceptor(c.logger),
	}
	streamInterceptors := []g.StreamServerInterceptor{
		interceptor.NewRecoverStreamServerInterceptor(c.logger),
	}

	serverOptions = append(serverOptions, grpc_middleware.WithUnaryServerChain(unaryInterceptors...))
	serverOptions = append(serverOptions, grpc_middleware.WithStreamServerChain(streamInterceptors...))

	// stats handlers
	statsHandlers := []s.Handler{
		stats.NewLoggerHandler(c.logger),
		stats.NewContextHandler(c.config),
	}
	if c.application.HasComponent(metrics.ComponentName) {
		statsHandlers = append(statsHandlers, stats.NewMetricHandler())
	}

	serverOptions = append(serverOptions, stats.WithStatsHandlerServerChain(statsHandlers...))

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

func (c *Component) GetServiceInfo() map[string]g.ServiceInfo {
	return c.server.GetServiceInfo()
}
