package system

import (
	"time"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
)

type SystemService struct {
	application *shadow.Application
	config      *config.Resource
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
	location, err := time.LoadLocation(s.config.GetString(ConfigSystemTimezone))
	if err != nil {
		return err
	}

	time.Local = location

	return nil
}
