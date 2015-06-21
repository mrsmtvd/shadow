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

	if r.config.GetBool("debug") {
		r.logger.Level = logrus.DebugLevel
	}

	return nil
}

func (r *Logger) Run() (err error) {
	r.Get(r.GetName()).Info("Logger start")

	return nil
}

func (r *Logger) Get(key string) *logrus.Entry {
	return r.logger.WithField("component", key)
}
