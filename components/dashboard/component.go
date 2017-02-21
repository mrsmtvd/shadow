package dashboard

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
)

type Component struct {
	application shadow.Application
	config      *config.Component
	logger      logger.Logger
	router      *Router
	renderer    *Renderer
}

func (c *Component) GetName() string {
	return "dashboard"
}

func (c *Component) GetVersion() string {
	return "1.0.0"
}

func (c *Component) Init(a shadow.Application) (err error) {
	resourceConfig, err := a.GetComponent("config")
	if err != nil {
		return err
	}
	c.config = resourceConfig.(*config.Component)

	c.application = a

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) error {
	c.logger = logger.NewOrNop(c.GetName(), c.application)

	if err := c.loadTemplates(); err != nil {
		return err
	}

	c.loadMenu()
	c.loadRoutes()

	go func(router *Router) {
		defer wg.Done()

		http.HandleFunc("/", func(out http.ResponseWriter, in *http.Request) {
			router.ServeHTTP(out, in)
		})

		// TODO: ssl
		addr := fmt.Sprintf("%s:%d", c.config.GetString(ConfigDashboardHost), c.config.GetInt(ConfigDashboardPort))

		c.logger.Info("Running service", map[string]interface{}{
			"addr": addr,
			"pid":  os.Getpid(),
		})

		if err := http.ListenAndServe(addr, c.router); err != nil {
			c.logger.Fatalf("Could not start frontend [%d]: %s\n", os.Getpid(), err.Error())
		}
	}(c.router)

	return nil
}
