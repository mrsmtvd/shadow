package logger

import (
	"flag"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
	"github.com/rs/xlog"
)

const (
	FieldAppName    = "app-name"
	FieldAppVersion = "app-version"
	FieldAppBuild   = "app-build"
	FieldComponent  = "component"
	FieldHostname   = "hostname"
)

type Resource struct {
	application *shadow.Application

	config  *config.Resource
	loggers map[string]Logger

	mutex        sync.RWMutex
	loggerConfig xlog.Config
}

func (r *Resource) GetName() string {
	return "logger"
}

func (r *Resource) Init(a *shadow.Application) error {
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	r.config = resourceConfig.(*config.Resource)

	r.application = a

	r.loggers = make(map[string]Logger, 1)

	return nil
}

func (r *Resource) Run() error {
	r.loggerConfig = xlog.Config{
		Output: xlog.NewConsoleOutput(),
		Level:  r.getLevel(),
		Fields: r.getDefaultFields(),
	}

	log.SetOutput(r.Get(r.GetName()))

	return nil
}

func (r *Resource) logConfig() {
	globalConfig := r.config.GetGlobalConf()
	fields := map[string]interface{}{
		"config.prefix": globalConfig.EnvPrefix,
		"config.file":   globalConfig.Filename,
	}

	for key := range r.config.GetAllValues() {
		fields[key] = r.config.Get(key)
	}

	logger := r.Get("config")
	logger.Info("Init config", fields)

	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == config.FlagConfig && f.Value.String() != "" {
			logger.Infof("Use config from %s", f.Value.String())
		}
	})
}

func (r *Resource) Get(key string) Logger {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r, ok := r.loggers[key]; ok {
		return r
	}

	l := newLogger(r.loggerConfig)
	l.SetField(FieldComponent, key)

	r.loggers[key] = l

	return l
}

func (r *Resource) getLevel() xlog.Level {
	switch r.config.GetInt(ConfigLoggerLevel) {
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

func (r *Resource) getDefaultFields() map[string]interface{} {
	fields := map[string]interface{}{
		FieldAppName:    r.application.Name,
		FieldAppVersion: r.application.Version,
		FieldAppBuild:   r.application.Build,
	}

	if hostname, err := os.Hostname(); err == nil {
		fields[FieldHostname] = hostname
	}

	fieldsFromConfig := r.config.GetString(ConfigLoggerFields)
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
