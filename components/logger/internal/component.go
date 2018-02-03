package internal

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
	"github.com/rs/xlog"
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

func (c *Component) GetName() string {
	return logger.ComponentName
}

func (c *Component) GetVersion() string {
	return logger.ComponentVersion
}

func (c *Component) GetDependencies() []shadow.Dependency {
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
	log.SetOutput(c.Get(c.GetName()))

	return nil
}

func (c *Component) Get(key string) logger.Logger {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if r, ok := c.loggers[key]; ok {
		return r
	}

	loggerConfig := xlog.Config{
		Output: xlog.NewConsoleOutput(),
		Level:  c.getXLogLevel(),
		Fields: c.getFields(),
	}

	loggerConfig.Fields[fieldComponent] = key

	l := newLogger(loggerConfig)
	c.loggers[key] = l

	return l
}

func (c *Component) getXLogLevel() xlog.Level {
	switch c.config.IntDefault(logger.ConfigLevel, 5) {
	case logger.LevelEmergency:
		return xlog.LevelFatal
	case logger.LevelAlert:
		return xlog.LevelFatal
	case logger.LevelCritical:
		return xlog.LevelFatal
	case logger.LevelError:
		return xlog.LevelError
	case logger.LevelWarning:
		return xlog.LevelWarn
	case logger.LevelNotice:
		return xlog.LevelInfo
	case logger.LevelInformational:
		return xlog.LevelInfo
	case logger.LevelDebug:
		return xlog.LevelDebug
	}

	return xlog.LevelInfo
}

func (c *Component) getFields() map[string]interface{} {
	fields := c.parseFields(c.config.String(logger.ConfigFields))

	if _, ok := fields[fieldComponent]; ok {
		delete(fields, fieldComponent)
	}

	fields[fieldAppName] = c.application.GetName()
	fields[fieldAppVersion] = c.application.GetVersion()
	fields[fieldAppBuild] = c.application.GetBuild()

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
