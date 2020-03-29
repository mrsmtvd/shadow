package internal

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/alexedwards/scs"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logging"
)

type Component struct {
	mutex sync.RWMutex

	application    shadow.Application
	components     []shadow.Component
	config         config.Component
	logger         logging.Logger
	renderer       *Renderer
	sessionManager *scs.Manager
	router         *Router
	server         *http.Server

	registryAssetFS *sync.Map
}

func (c *Component) Name() string {
	return dashboard.ComponentName
}

func (c *Component) Version() string {
	return dashboard.ComponentVersion
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

func (c *Component) Init(a shadow.Application) (err error) {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.renderer = NewRenderer()
	c.registryAssetFS = new(sync.Map)

	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) (err error) {
	if c.components, err = a.GetComponents(); err != nil {
		return err
	}

	c.logger = logging.DefaultLazyLogger(c.Name())

	<-a.ReadyComponent(config.ComponentName)

	c.router = NewRouter(c.logger, c.config.Int(dashboard.ConfigPanicHandlerCallerSkip))

	c.initAssetFS()

	if err := c.initTemplates(); err != nil {
		return err
	}

	c.initMenu()
	c.initSession()

	if err := c.initAuth(); err != nil {
		return err
	}

	if err = c.initServeMux(); err != nil {
		return err
	}

	addr := net.JoinHostPort(c.config.String(dashboard.ConfigHost), c.config.String(dashboard.ConfigPort))
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		return fmt.Errorf("failed to listen [%d]: %s", os.Getpid(), err.Error())
	}

	c.logger.Info("Running service", "addr", addr, "pid", os.Getpid())

	srv := &http.Server{
		Handler: c.router,
	}

	c.mutex.Lock()
	c.server = srv
	c.mutex.Unlock()

	ready <- struct{}{}

	if err := srv.Serve(lis); err != nil && err != http.ErrServerClosed {
		c.logger.Errorf("Failed to serve [%d]: %s\n", os.Getpid(), err.Error())
		return err
	}

	return nil
}

func (c *Component) Shutdown() error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.server != nil {
		return c.server.Shutdown(context.Background())
	}

	return nil
}

func (c *Component) Renderer() dashboard.Renderer {
	return c.renderer
}
