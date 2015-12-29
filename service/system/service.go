package system

import (
	"time"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
)

type SystemService struct {
	Application *shadow.Application
	Config      *resource.Config
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
	s.Config = resourceConfig.(*resource.Config)

	location, err := time.LoadLocation(s.Config.GetString("system.timezone"))
	if err != nil {
		return err
	}

	time.Local = location

	return nil
}
