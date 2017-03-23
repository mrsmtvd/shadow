package logger

import (
	"flag"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/rs/xlog"
)

const (
	ComponentName = "logger"

	FieldAppName    = "app-name"
	FieldAppVersion = "app-version"
	FieldAppBuild   = "app-build"
	FieldComponent  = "component"
	FieldHostname   = "hostname"
)

type Component struct {
	application shadow.Application

	config  *config.Component
	loggers map[string]Logger

	mutex        sync.RWMutex
	loggerConfig xlog.Config
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
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(*config.Component)
	c.loggers = make(map[string]Logger, 1)

	return nil
}

func (c *Component) Run() error {
	c.loggerConfig = xlog.Config{
		Output: xlog.NewConsoleOutput(),
		Level:  c.getLevel(),
		Fields: c.getDefaultFields(),
	}

	log.SetOutput(c.Get(c.GetName()))

	return nil
}

func (c *Component) logConfig() {
	globalConfig := c.config.GetGlobalConf()
	fields := map[string]interface{}{
		"config.prefix": globalConfig.EnvPrefix,
		"config.file":   globalConfig.Filename,
	}

	for key := range c.config.GetAllValues() {
		fields[key] = c.config.Get(key)
	}

	logger := c.Get("config")
	logger.Info("Init config", fields)

	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == config.FlagConfig && f.Value.String() != "" {
			logger.Infof("Use config from %s", f.Value.String())
		}
	})
}

func (c *Component) Get(key string) Logger {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if r, ok := c.loggers[key]; ok {
		return r
	}

	l := newLogger(c.loggerConfig)
	l.SetField(FieldComponent, key)

	c.loggers[key] = l

	return l
}

func (c *Component) getLevel() xlog.Level {
	switch c.config.GetInt(ConfigLoggerLevel) {
	case 0:
		return xlog.LevelFatal
	case 1:
		return xlog.LevelFatal
	case 2:
		return xlog.LevelFatal
	case 3:
		return xlog.LevelError
	case 4:
		return xlog.LevelWarn
	case 5:
		return xlog.LevelInfo
	case 6:
		return xlog.LevelInfo
	case 7:
		return xlog.LevelDebug
	}

	return xlog.LevelInfo
}

func (c *Component) getDefaultFields() map[string]interface{} {
	fields := map[string]interface{}{
		FieldAppName:    c.application.GetName(),
		FieldAppVersion: c.application.GetVersion(),
		FieldAppBuild:   c.application.GetBuild(),
	}

	if hostname, err := os.Hostname(); err == nil {
		fields[FieldHostname] = hostname
	}

	fieldsFromConfig := c.config.GetString(ConfigLoggerFields)
	if len(fieldsFromConfig) > 0 {
		var parts []string

		for _, tag := range strings.Split(fieldsFromConfig, ",") {
			parts = strings.Split(tag, "=")

			if len(parts) > 1 {
				fields[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	return fields
}
