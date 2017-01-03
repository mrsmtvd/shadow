package shadow // import "github.com/kihamo/shadow"

import (
	"fmt"
	"sync"
)

//go:generate goimports -w ./
//go:generate sh -c "cd components/alerts && go-bindata-assetfs -pkg=alerts templates/..."
//go:generate sh -c "cd components/dashboard && go-bindata-assetfs -pkg=dashboard templates/... public/..."
//go:generate sh -c "cd components/mail && go-bindata-assetfs -pkg=mail templates/..."
//go:generate sh -c "cd components/workers && go-bindata-assetfs -pkg=workers templates/..."

type Component interface {
	GetName() string
	GetVersion() string
	Init(*Application) error
}

type ComponentRunner interface {
	Run() error
}

type ComponentAsyncRunner interface {
	Run(*sync.WaitGroup) error
}

type Application struct {
	components []Component

	Name    string
	Version string
	Build   string

	wg *sync.WaitGroup
}

func NewApplication(components []Component, name string, version string, build string) (*Application, error) {
	application := &Application{
		components: []Component{},
		Name:       name,
		Version:    version,
		Build:      build,
		wg:         new(sync.WaitGroup),
	}

	for i := range components {
		if err := application.RegisterComponent(components[i]); err != nil {
			return nil, err
		}
	}

	return application, nil
}

func (a *Application) Run() (err error) {
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

func (a *Application) GetComponent(n string) (Component, error) {
	for i := range a.components {
		if a.components[i].GetName() == n {
			return a.components[i], nil
		}
	}

	return nil, fmt.Errorf("Component \"%s\" not found", n)
}

func (a *Application) GetComponents() []Component {
	return a.components
}

func (a *Application) HasComponent(n string) bool {
	_, err := a.GetComponent(n)
	return err == nil
}

func (a *Application) RegisterComponent(c Component) error {
	if _, err := a.GetComponent(c.GetName()); err == nil {
		return fmt.Errorf("Component \"%s\" already exists", c.GetName())
	}

	a.components = append(a.components, c)
	return nil
}
