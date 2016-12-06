package logger

import (
	"flag"
	"log"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
	"github.com/rs/xlog"
)

type Logger struct {
	config       *config.Config
	logger       xlog.Logger
	loggerConfig xlog.Config
}

func (r *Logger) GetName() string {
	return "logger"
}

func (r *Logger) GetConfigVariables() []config.ConfigVariable {
	return []config.ConfigVariable{
		config.ConfigVariable{
			Key:   "logger.level",
			Value: 5,
			Usage: "Log level",
		},
	}
}

func (r *Logger) Init(a *shadow.Application) error {
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	r.config = resourceConfig.(*config.Config)

	r.loggerConfig = xlog.Config{
		Output: xlog.NewConsoleOutput(),
	}

	r.initLogger()

	return nil
}

func (r *Logger) Run() (err error) {
	var level xlog.Level

	if r.config.GetBool("debug") {
		level = xlog.LevelDebug
	} else {
		switch r.config.GetInt("logger.level") {
		case 1:
			level = xlog.LevelFatal
		case 2:
			level = xlog.LevelFatal
		case 3:
			level = xlog.LevelError
		case 4:
			level = xlog.LevelWarn
		case 5:
			level = xlog.LevelInfo
		case 6:
			level = xlog.LevelDebug
		}
	}

	if level != r.loggerConfig.Level {
		r.loggerConfig.Level = level
		r.initLogger()
	}

	r.logConfig()

	return nil
}

func (r *Logger) initLogger() {
	r.logger = xlog.New(r.loggerConfig)
	log.SetOutput(r.logger)
}

func (r *Logger) logConfig() {
	globalConfig := r.config.GetGlobalConf()
	fields := xlog.F{
		"config.prefix": globalConfig.EnvPrefix,
		"config.file":   globalConfig.Filename,
	}

	for key := range r.config.GetAll() {
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

func (r *Logger) Get(key string) xlog.Logger {
	logger := xlog.Copy(r.logger)
	logger.SetField("component", key)

	return logger
}
