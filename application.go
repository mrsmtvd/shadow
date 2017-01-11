package shadow

import (
	"fmt"
	"sync"
)

type Application interface {
	Run() error
	GetComponent(string) (Component, error)
	GetComponents() []Component
	HasComponent(string) bool
	RegisterComponent(Component) error
	GetName() string
	GetVersion() string
	GetBuild() string
}

type Component interface {
	GetName() string
	GetVersion() string
	Init(Application) error
}

type ComponentRunner interface {
	Run() error
}

type ComponentAsyncRunner interface {
	Run(*sync.WaitGroup) error
}

type App struct {
	components []Component

	name    string
	version string
	build   string

	wg *sync.WaitGroup
}

func NewApp(components []Component, name string, version string, build string) (Application, error) {
	application := &App{
		components: []Component{},
		name:       name,
		version:    version,
		build:      build,
		wg:         new(sync.WaitGroup),
	}

	for i := range components {
		if err := application.RegisterComponent(components[i]); err != nil {
			return nil, err
		}
	}

	return application, nil
}

func (a *App) Run() (err error) {
	for i := range a.components {
		if err = a.components[i].Init(a); err != nil {
			return err
		}
	}

	for i := range a.components {
		if runner, ok := a.components[i].(ComponentAsyncRunner); ok {
			a.wg.Add(1)
			if err := runner.Run(a.wg); err != nil {
				a.wg.Done()
				return err
			}
		} else if runner, ok := a.components[i].(ComponentRunner); ok {
			if err := runner.Run(); err != nil {
				return err
			}
		}
	}

	a.wg.Wait()
	return nil
}

func (a *App) GetComponent(n string) (Component, error) {
	for i := range a.components {
		if a.components[i].GetName() == n {
			return a.components[i], nil
		}
	}

	return nil, fmt.Errorf("Component \"%s\" not found", n)
}

func (a *App) GetComponents() []Component {
	return a.components
}

func (a *App) HasComponent(n string) bool {
	_, err := a.GetComponent(n)
	return err == nil
}

func (a *App) RegisterComponent(c Component) error {
	if _, err := a.GetComponent(c.GetName()); err == nil {
		return fmt.Errorf("Component \"%s\" already exists", c.GetName())
	}

	a.components = append(a.components, c)
	return nil
}

func (a *App) GetName() string {
	return a.name
}

func (a *App) GetVersion() string {
	return a.version
}

func (a *App) GetBuild() string {
	return a.build
}
