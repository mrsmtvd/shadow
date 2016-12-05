package logger

import (
	"flag"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
)

type Logger struct {
	config *config.Config
	logger *logrus.Logger
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

	r.logger = logrus.StandardLogger()
	formatter := r.logger.Formatter
	if textFormatter, ok := formatter.(*logrus.TextFormatter); ok {
		textFormatter.FullTimestamp = true
		textFormatter.TimestampFormat = time.RFC3339Nano
	}

	return nil
}

func (r *Logger) Run() (err error) {
	if r.config.GetBool("debug") {
		r.logger.Level = logrus.DebugLevel
	} else {
		switch r.config.GetInt("logger.level") {
		case 1:
			r.logger.Level = logrus.PanicLevel
		case 2:
			r.logger.Level = logrus.FatalLevel
		case 3:
			r.logger.Level = logrus.ErrorLevel
		case 4:
			r.logger.Level = logrus.WarnLevel
		case 5:
			r.logger.Level = logrus.InfoLevel
		case 6:
			r.logger.Level = logrus.DebugLevel
		}
	}

	fields := logrus.Fields{}
	for key := range r.config.GetAll() {
		fields[key] = r.config.Get(key)
	}

	logger := r.Get("config").WithFields(fields)
	logger.WithFields(logrus.Fields{
		"prefix": r.config.GetGlobalConf().EnvPrefix,
		"file":   r.config.GetGlobalConf().Filename,
	}).Info("Init config")

	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == config.FlagConfig && f.Value.String() != "" {
			logger.Infof("Use config from %s", f.Value.String())
		}
	})

	return nil
}

func (r *Logger) Get(key string) *logrus.Entry {
	return r.logger.WithField("component", key)
}
