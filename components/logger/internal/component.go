package internal

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/logger/output"
)

const (
	fieldAppName    = "app-name"
	fieldAppVersion = "app-version"
	fieldAppBuild   = "app-build"
	fieldComponent  = "component"
	fieldHostname   = "hostname"
)

type Component struct {
	application shadow.Application
	config      config.Component
	loggers     map[string]logger.Logger
	mutex       sync.RWMutex
}

func (c *Component) Name() string {
	return logger.ComponentName
}

func (c *Component) Version() string {
	return logger.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.loggers = make(map[string]logger.Logger, 1)

	return nil
}

func (c *Component) Run() error {
	log.SetOutput(c.Get(c.Name()))

	return nil
}

func (c *Component) Get(key string) logger.Logger {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if r, ok := c.loggers[key]; ok {
		return r
	}

	fields := c.getFields()
	fields[fieldComponent] = key
	c.loggers[key] = output.NewConsoleOutput(c.getLevel(), fields)
	return c.loggers[key]
}

func (c *Component) getLevel() logger.Level {
	return logger.Level(c.config.IntDefault(logger.ConfigLevel, 5))
}

func (c *Component) getFields() map[string]interface{} {
	fields := c.parseFields(c.config.String(logger.ConfigFields))

	if _, ok := fields[fieldComponent]; ok {
		delete(fields, fieldComponent)
	}

	fields[fieldAppName] = c.application.Name()
	fields[fieldAppVersion] = c.application.Version()
	fields[fieldAppBuild] = c.application.Build()

	if hostname, err := os.Hostname(); err == nil {
		fields[fieldHostname] = hostname
	}

	return fields
}

func (c *Component) parseFields(f string) map[string]interface{} {
	fields := map[string]interface{}{}

	if len(f) == 0 {
		return fields
	}

	var parts []string

	for _, tag := range strings.Split(f, ",") {
		parts = strings.Split(tag, "=")

		if len(parts) > 1 {
			fields[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return fields
}
