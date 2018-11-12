package shadow

import (
	"errors"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/bugsnag/osext"
	"github.com/deckarep/golang-set"
	"golang.org/x/sync/errgroup"
)

var (
	startTime = time.Now().UTC()

	DefaultApplication = NewApp()
)

type Application interface {
	Run() error
	GetComponent(string) Component
	GetComponents() ([]Component, error)
	HasComponent(string) bool
	RegisterComponent(Component) error
	Name() string
	Version() string
	Build() string
	BuildDate() *time.Time
	StartDate() *time.Time
	Uptime() time.Duration
}

type Component interface {
	Name() string
	Version() string
}

type ComponentInit interface {
	Init(Application) error
}

type ComponentRunner interface {
	Run() error
}

type ComponentDependency interface {
	Dependencies() []Dependency
}

type ComponentShutdown interface {
	Shutdown() error
}

type Dependency struct {
	Name     string
	Required bool
}

type App struct {
	_ [4]byte // atomic requires 64-bit alignment for struct field access

	components        map[string]Component
	resolveComponents []Component

	name    string
	version string
	build   string

	resolved bool
	shutdown int64

	runDone        chan error
	shutdownSignal chan os.Signal
	closers        []func() error
}

var (
	buildDate *time.Time
)

func init() {
	if b, err := osext.Executable(); err == nil {
		if f, err := os.Stat(b); err == nil {
			d := f.ModTime().UTC()
			buildDate = &d
		}
	}
}

func NewApp() *App {
	application := &App{
		components: map[string]Component{},
	}

	return application
}

func (a *App) Run() (err error) {
	if a.shutdownSignal != nil {
		return errors.New("already running")
	}

	components, err := a.GetComponents()
	if err != nil {
		return err
	}

	a.runDone = make(chan error)

	// listen os signal
	a.shutdownSignal = make(chan os.Signal)
	signal.Notify(a.shutdownSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	a.closers = make([]func() error, 0, len(components))

	for i := range components {
		if closer, ok := components[i].(ComponentShutdown); ok {
			a.closers = append([]func() error{closer.Shutdown}, a.closers...)
		}

		if init, ok := components[i].(ComponentInit); ok {
			if err = init.Init(a); err != nil {
				return err
			}
		}
	}

	if len(a.closers) > 0 {
		go func() {
			<-a.shutdownSignal

			atomic.StoreInt64(&a.shutdown, 1)

			signal.Stop(a.shutdownSignal)
			close(a.shutdownSignal)

			var shutdownEG errgroup.Group

			for _, closer := range a.closers {
				shutdownEG.Go(closer)
			}

			a.runDone <- shutdownEG.Wait()
		}()
	}

	errCh := make(chan error)

	go func() {
		select {
		case err := <-errCh:
			if atomic.LoadInt64(&a.shutdown) == 0 {
				a.closers = append([]func() error{
					func() error {
						return err
					},
				}, a.closers...)

				a.shutdownSignal <- syscall.SIGQUIT
			}

			close(errCh)
			return
		}
	}()

	for i := range components {
		if runner, ok := components[i].(ComponentRunner); ok {
			go func() {
				if err := runner.Run(); err != nil {
					errCh <- err
				}
			}()
		}
	}

	return <-a.runDone
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
	if a.HasComponent(c.Name()) {
		return errors.New("component \"" + c.Name() + "\" already exists")
	}

	a.components[c.Name()] = c
	a.resolved = false
	return nil
}

func (a *App) MustRegisterComponent(c Component) {
	if err := a.RegisterComponent(c); err != nil {
		panic(err)
	}
}

func (a *App) SetName(name string) {
	a.name = name
}

func (a *App) Name() string {
	return a.name
}

func (a *App) SetVersion(version string) {
	a.version = version
}

func (a *App) Version() string {
	return a.version
}

func (a *App) SetBuild(build string) {
	a.build = build
}

func (a *App) Build() string {
	return a.build
}

func (a *App) BuildDate() *time.Time {
	return buildDate
}

func (a *App) StartDate() *time.Time {
	return &startTime
}

func (a *App) Uptime() time.Duration {
	return time.Now().UTC().Sub(startTime)
}

func (a *App) resolveDependencies() error {
	a.resolveComponents = make([]Component, 0, len(a.components))

	cmpDependencies := make(map[string]mapset.Set)
	for _, cmp := range a.components {
		dependencySet := mapset.NewSet()

		if cmpDependency, ok := cmp.(ComponentDependency); ok {
			for _, dep := range cmpDependency.Dependencies() {
				if dep.Required {
					if !a.HasComponent(dep.Name) {
						return errors.New("Component \"" + cmp.Name() + "\" has required dependency \"" + dep.Name + "\"")
					}
				} else if !a.HasComponent(dep.Name) {
					cmpDependencies[dep.Name] = mapset.NewSet()
				}

				dependencySet.Add(dep.Name)
			}
		}

		cmpDependencies[cmp.Name()] = dependencySet
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

func SetName(name string) {
	DefaultApplication.SetName(name)
}

func SetVersion(version string) {
	DefaultApplication.SetVersion(version)
}

func SetBuild(build string) {
	DefaultApplication.SetBuild(build)
}

func MustRegisterComponent(c Component) {
	DefaultApplication.MustRegisterComponent(c)
}

func RegisterComponent(c Component) error {
	return DefaultApplication.RegisterComponent(c)
}

func Run() error {
	return DefaultApplication.Run()
}
