package shadow // import "github.com/kihamo/shadow"

import (
	"errors"
	"sync"
)

//go:generate goimports -w ./
//go:generate sh -c "cd service/frontend && go-bindata-assetfs -pkg=frontend templates/... public/..."
//go:generate sh -c "cd service/system && go-bindata-assetfs -pkg=system templates/..."

type ContextItem interface {
	GetName() string
	Init(*Application) error
}

type ContextItemRunner interface {
	Run() error
}

type ContextItemAsyncRunner interface {
	Run(*sync.WaitGroup) error
}

type Resource interface {
	ContextItem
}

type Service interface {
	ContextItem
}

type Application struct {
	resources []Resource
	services  []Service

	Name    string
	Version string
	Build   string

	wg *sync.WaitGroup
}

func NewApplication(resources []Resource, services []Service, name string, version string, build string) (*Application, error) {
	application := &Application{
		resources: []Resource{},
		services:  []Service{},
		Name:      name,
		Version:   version,
		Build:     build,
		wg:        new(sync.WaitGroup),
	}

	for i := range resources {
		if err := application.RegisterResource(resources[i]); err != nil {
			return nil, err
		}
	}

	for i := range services {
		if err := application.RegisterService(services[i]); err != nil {
			return nil, err
		}
	}

	return application, nil
}

func (a *Application) Run() (err error) {
	// Resources
	resources := a.GetResources()

	for i := range resources {
		if err = resources[i].Init(a); err != nil {
			return err
		}
	}

	for i := range resources {
		if err = a.run(resources[i]); err != nil {
			return err
		}
	}

	// Services
	services := a.GetServices()

	for i := range services {
		if err = services[i].Init(a); err != nil {
			return err
		}
	}

	for i := range services {
		if err = a.run(services[i]); err != nil {
			return err
		}
	}

	a.wg.Wait()
	return nil
}

func (a *Application) run(item ContextItem) error {
	if runner, ok := item.(ContextItemAsyncRunner); ok {
		a.wg.Add(1)
		if err := runner.Run(a.wg); err != nil {
			a.wg.Done()
			return err
		}
	} else if runner, ok := item.(ContextItemRunner); ok {
		if err := runner.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (a *Application) GetResource(name string) (Resource, error) {
	for i := range a.resources {
		if a.resources[i].GetName() == name {
			return a.resources[i], nil
		}
	}

	return nil, errors.New("Resource \"" + name + "\" not found")
}

func (a *Application) GetResources() []Resource {
	return a.resources
}

func (a *Application) HasResource(name string) bool {
	_, err := a.GetResource(name)
	return err == nil
}

func (a *Application) RegisterResource(resource Resource) error {
	if _, err := a.GetResource(resource.GetName()); err == nil {
		return errors.New("Resource \"" + resource.GetName() + "\" already exists")
	}

	a.resources = append(a.resources, resource)
	return nil
}

func (a *Application) GetService(name string) (Service, error) {
	for i := range a.services {
		if a.services[i].GetName() == name {
			return a.services[i], nil
		}
	}

	return nil, errors.New("Service \"" + name + "\" not found")
}

func (a *Application) GetServices() []Service {
	return a.services
}

func (a *Application) HasService(name string) bool {
	_, err := a.GetService(name)
	return err == nil
}

func (a *Application) RegisterService(service Service) error {
	if _, err := a.GetService(service.GetName()); err == nil {
		return errors.New("Service \"" + service.GetName() + "\" already exists")
	}

	a.services = append(a.services, service)
	return nil
}
