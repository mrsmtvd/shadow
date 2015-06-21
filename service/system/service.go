package system

import (
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

	return nil
}
