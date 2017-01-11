package dashboard

import (
	"encoding/base64"
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

	mutex     sync.RWMutex
	renderer  *Renderer
	authToken string
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

	c.generateAuthToken(c.config.GetString(ConfigFrontendAuthUser), c.config.GetString(ConfigFrontendAuthPassword))

	go func(router *Router) {
		defer wg.Done()

		http.HandleFunc("/", func(out http.ResponseWriter, in *http.Request) {
			router.ServeHTTP(out, in)
		})

		// TODO: ssl
		addr := fmt.Sprintf("%s:%d", c.config.GetString(ConfigFrontendHost), c.config.GetInt(ConfigFrontendPort))

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

func (c *Component) generateAuthToken(user, password string) {
	token := ""
	if user != "" {
		token = "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password))
	}

	c.mutex.Lock()
	c.authToken = token
	c.mutex.Unlock()
}
