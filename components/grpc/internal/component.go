package internal

import (
	"net"
	"os"
	"sync"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/grpc"
	"github.com/kihamo/shadow/components/grpc/server"
	"github.com/kihamo/shadow/components/grpc/stats"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logging"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	s "google.golang.org/grpc/stats"
)

type Component struct {
	mutex sync.RWMutex

	application shadow.Application
	config      config.Component
	logger      logging.Logger
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
			Name: logging.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) error {
	c.logger = logging.DefaultLogger().Named(c.Name())
	grpclog.SetLoggerV2(grpc.NewLogger(c.logger))

	<-a.ReadyComponent(config.ComponentName)

	unaryInterceptors := make([]g.UnaryServerInterceptor, 0)
	streamInterceptors := make([]g.StreamServerInterceptor, 0)
	statsHandlers := []s.Handler{
		stats.NewContextHandler(c.config),
	}

	components, err := a.GetComponents()
	if err != nil {
		return err
	}

	for _, cmp := range components {
		if cmpUnaryServerInterceptors, ok := cmp.(grpc.HasUnaryServerInterceptors); ok {
			unaryInterceptors = append(unaryInterceptors, cmpUnaryServerInterceptors.GrpcUnaryServerInterceptors()...)
		}

		if cmpStreamServerInterceptors, ok := cmp.(grpc.HasStreamServerInterceptors); ok {
			streamInterceptors = append(streamInterceptors, cmpStreamServerInterceptors.GrpcStreamServerInterceptors()...)
		}

		if cmpStatsHandlers, ok := cmp.(grpc.HasStatsHandlers); ok {
			statsHandlers = append(statsHandlers, cmpStatsHandlers.GrpcStatsHandlers()...)
		}
	}

	srv := server.NewDefaultServerWithCustomOptions(unaryInterceptors, streamInterceptors, statsHandlers)

	for _, cmp := range components {
		if cmpGrpc, ok := cmp.(grpc.HasGrpcServer); ok {
			cmpGrpc.RegisterGrpcServer(srv)
		}
	}

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, healthServer)

	for service, info := range srv.GetServiceInfo() {
		healthServer.SetServingStatus(service, grpc_health_v1.HealthCheckResponse_SERVING)

		for _, method := range info.Methods {
			c.logger.Debug("Add method /" + service + "/" + method.Name)
		}
	}

	if c.config.Bool(grpc.ConfigReflectionEnabled) {
		reflection.Register(srv)
	}

	addr := net.JoinHostPort(c.config.String(grpc.ConfigHost), c.config.String(grpc.ConfigPort))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		c.logger.Errorf("Failed to listen [%d]: %s\n", os.Getpid(), err.Error())
		return err
	}

	c.logger.Info("Running service",
		"addr", addr,
		"pid", os.Getpid(),
	)

	c.mutex.Lock()
	c.server = srv
	c.mutex.Unlock()

	ready <- struct{}{}

	if err := srv.Serve(lis); err != nil {
		c.logger.Errorf("Failed to serve [%d]: %s\n", os.Getpid(), err.Error())
		return err
	}

	return nil
}

func (c *Component) Shutdown() error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.server != nil {
		c.server.GracefulStop()
	}

	return nil
}

func (c *Component) GetServiceInfo() map[string]g.ServiceInfo {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.server == nil {
		return nil
	}

	return c.server.GetServiceInfo()
}
