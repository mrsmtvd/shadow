package internal

import (
	"net"
	"sync"
	"time"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
)

type Component struct {
	application shadow.Application
	config      config.Component

	mutex  sync.RWMutex
	tracer opentracing.Tracer
}

func (c *Component) Name() string {
	return tracing.ComponentName
}

func (c *Component) Version() string {
	return tracing.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: logger.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	return nil
}

func (c *Component) Run() error {
	cfg := jconfig.Configuration{
		Disabled:    false,
		ServiceName: c.application.Name(),
		Sampler: &jconfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
	}

	sender, _ := jaeger.NewUDPTransport(net.JoinHostPort(c.config.String(tracing.ConfigAgentHost), c.config.String(tracing.ConfigAgentPort)), 0)

	options := []jconfig.Option{
		jconfig.Reporter(jaeger.NewRemoteReporter(sender, jaeger.ReporterOptions.BufferFlushInterval(1*time.Second))),
	}

	if c.application.HasComponent(logger.ComponentName) {
		log := logger.NewOrNop(c.Name(), c.application)
		options = append(options, jconfig.Logger(NewLogger(log)))
	}

	tracer, _, err := cfg.NewTracer(options...)
	if err != nil {
		return err
	}

	// if is global?
	c.mutex.Lock()
	c.tracer = tracer
	c.mutex.Unlock()

	//opentracing.SetGlobalTracer(tracer)

	return nil
}

func (c *Component) Tracer() opentracing.Tracer {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.tracer
}
