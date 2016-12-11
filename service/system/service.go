package system

import (
	"time"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
	"github.com/kihamo/shadow/resource/logger"
)

type SystemService struct {
	application *shadow.Application
	config      *config.Resource
	logger      logger.Logger
}

func (s *SystemService) GetName() string {
	return "system"
}

func (s *SystemService) Init(a *shadow.Application) error {
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	s.config = resourceConfig.(*config.Resource)

	s.application = a

	return nil
}

func (s *SystemService) Run() error {
	location, err := time.LoadLocation(s.config.GetString("system.timezone"))
	if err != nil {
		return err
	}

	time.Local = location

	if resourceLogger, err := s.application.GetResource("logger"); err == nil {
		s.logger = resourceLogger.(*logger.Resource).Get(s.GetName())
	} else {
		s.logger = logger.NopLogger
	}

	return nil
}
