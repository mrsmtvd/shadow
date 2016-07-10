package resource

import (
	"github.com/Sirupsen/logrus"
	"github.com/kihamo/shadow"
)

type Logger struct {
	config *Config
	logger *logrus.Logger
}

func (r *Logger) GetName() string {
	return "logger"
}

func (r *Logger) GetConfigVariables() []ConfigVariable {
	return []ConfigVariable{
		ConfigVariable{
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
	r.config = resourceConfig.(*Config)

	r.logger = logrus.StandardLogger()
	formatter := r.logger.Formatter
	if textFormatter, ok := formatter.(*logrus.TextFormatter); ok {
		textFormatter.FullTimestamp = true
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

	r.Get(r.GetName()).Info("Logger start")

	return nil
}

func (r *Logger) Get(key string) *logrus.Entry {
	return r.logger.WithField("component", key)
}
