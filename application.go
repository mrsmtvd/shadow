package shadow

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bugsnag/osext"
	"github.com/deckarep/golang-set"
)

type Application interface {
	Run() error
	GetComponent(string) Component
	GetComponents() ([]Component, error)
	HasComponent(string) bool
	RegisterComponent(Component) error
	GetName() string
	GetVersion() string
	GetBuild() string
	GetBuildDate() *time.Time
}

type Component interface {
	GetName() string
	GetVersion() string
}

type ComponentInit interface {
	Init(Application) error
}

type ComponentRunner interface {
	Run() error
}

type ComponentDependency interface {
	GetDependencies() []Dependency
}

type ComponentAsyncRunner interface {
	Run(*sync.WaitGroup) error
}

type Dependency struct {
	Name     string
	Required bool
}

type App struct {
	components        map[string]Component
	resolveComponents []Component

	name    string
	version string
	build   string

	wg       *sync.WaitGroup
	run      bool
	resolved bool
}

var buildDate *time.Time

func init() {
	if b, err := osext.Executable(); err == nil {
		if f, err := os.Stat(b); err == nil {
			d := f.ModTime()
			buildDate = &d
		}
	}
}

func NewApp(name string, version string, build string, components []Component) (Application, error) {
	application := &App{
		components: map[string]Component{},
		name:       name,
		version:    version,
		build:      build,
		wg:         new(sync.WaitGroup),
		run:        false,
	}

	for i := range components {
		if err := application.RegisterComponent(components[i]); err != nil {
			return nil, err
		}
	}

	return application, nil
}

func (a *App) Run() (err error) {
	if a.run {
		return errors.New("Already running")
	}

	components, err := a.GetComponents()
	if err != nil {
		return err
	}

	a.run = true

	for i := range components {
		if init, ok := components[i].(ComponentInit); ok {
			if err = init.Init(a); err != nil {
				return err
			}
		}
	}

	for i := range components {
		if runner, ok := components[i].(ComponentAsyncRunner); ok {
			a.wg.Add(1)
			if err := runner.Run(a.wg); err != nil {
				return err
			}
		} else if runner, ok := components[i].(ComponentRunner); ok {
			if err := runner.Run(); err != nil {
				return err
			}
		}
	}

	a.wg.Wait()
	return nil
}

func (a *App) GetComponent(n string) Component {
	if cmp, ok := a.components[n]; ok {
		return cmp
	}

	return nil
}

func (a *App) GetComponents() ([]Component, error) {
	if !a.resolved {
		if err := a.resolveDependencies(); err != nil {
			return nil, err
		}
	}

	return a.resolveComponents, nil
}

func (a *App) HasComponent(n string) bool {
	return a.GetComponent(n) != nil
}

func (a *App) RegisterComponent(c Component) error {
	if a.HasComponent(c.GetName()) {
		return fmt.Errorf("Component \"%s\" already exists", c.GetName())
	}

	a.components[c.GetName()] = c
	a.resolved = false
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

func (a *App) GetBuildDate() *time.Time {
	return buildDate
}

func (a *App) resolveDependencies() error {
	a.resolveComponents = make([]Component, 0, len(a.components))

	cmpDependencies := make(map[string]mapset.Set)
	for _, cmp := range a.components {
		dependencySet := mapset.NewSet()

		if cmpDependency, ok := cmp.(ComponentDependency); ok {
			for _, dep := range cmpDependency.GetDependencies() {
				if dep.Required {
					if !a.HasComponent(dep.Name) {
						return fmt.Errorf("Component \"%s\" has required dependency \"%s\"", cmp.GetName(), dep.Name)
					}
				} else if !a.HasComponent(dep.Name) {
					cmpDependencies[dep.Name] = mapset.NewSet()
				}

				dependencySet.Add(dep.Name)
			}
		}

		cmpDependencies[cmp.GetName()] = dependencySet
	}

	var components []string

	for len(cmpDependencies) > 0 {
		readySet := mapset.NewSet()
		for name, deps := range cmpDependencies {
			if deps.Cardinality() == 0 {
				readySet.Add(name)
			}
		}

		if readySet.Cardinality() == 0 {
			return errors.New("Circular dependency found")
		}

		for name := range readySet.Iter() {
			delete(cmpDependencies, name.(string))
			components = append(components, name.(string))
		}

		for name, deps := range cmpDependencies {
			diff := deps.Difference(readySet)
			cmpDependencies[name] = diff
		}
	}

	for _, name := range components {
		if cmp := a.GetComponent(name); cmp != nil {
			a.resolveComponents = append(a.resolveComponents, cmp)
		}
	}

	a.resolved = true
	return nil
}
