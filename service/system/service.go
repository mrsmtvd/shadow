package system

import (
	"time"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
)

type SystemService struct {
	Application *shadow.Application
}

func (s *SystemService) GetName() string {
	return "system"
}

func (s *SystemService) Init(a *shadow.Application) error {
	s.Application = a

	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}

	location, err := time.LoadLocation(resourceConfig.(*config.Resource).GetString("system.timezone"))
	if err != nil {
		return err
	}

	time.Local = location

	return nil
}
