package shadow

import (
	"errors"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/bugsnag/osext"
	"golang.org/x/sync/errgroup"
)

var (
	startTime = time.Now().UTC()

	DefaultApplication = NewApp()
)

type Application interface {
	Run() error
	Name() string
	Version() string
	Build() string
	BuildDate() *time.Time
	StartDate() *time.Time
	Uptime() time.Duration
	Shutdown() error

	GetComponent(string) Component
	GetComponents() ([]Component, error)
	HasComponent(string) bool
	RegisterComponent(Component) error

	StatusComponent(string) componentStatus
	WatchComponentStatus(status componentStatus, name string, names ...string) <-chan struct{}
	ReadyComponent(name string, names ...string) <-chan struct{}
	RunningComponent(name string, names ...string) <-chan struct{}
	ShutdownComponent(name string, names ...string) <-chan struct{}
}

type App struct {
	_ [4]byte // atomic requires 64-bit alignment for struct field access

	components *components
	running    int64
	shutdown   chan os.Signal

	name    string
	version string
	build   string
}

var (
	buildDate *time.Time
	sig       = []os.Signal{
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	}
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
		components: &components{},
		shutdown:   make(chan os.Signal, 1),
	}

	return application
}

func (a *App) Run() (err error) {
	if atomic.LoadInt64(&a.running) == 1 {
		return errors.New("already running")
	}

	atomic.StoreInt64(&a.running, 1)
	defer atomic.StoreInt64(&a.running, 0)

	components, err := a.components.All()
	if err != nil {
		return err
	}

	total := len(components)
	if total == 0 {
		return
	}

	// инициализируем компоненты
	for _, cmp := range components {
		if in, ok := cmp.instance.(ComponentInit); ok {
			if err := in.Init(a); err != nil {
				return err
			}
		}
	}

	// региструем shutdown функции
	closers := make([]func() error, 0, total)

	for _, cmp := range components {
		if cmp.closer == nil {
			continue
		}

		fn := func(component *component) func() error {
			return func() error {
				for _, dep := range component.ReverseDep() {
					<-a.ShutdownComponent(dep)
				}

				return component.Shutdown()
			}
		}(cmp)

		closers = append(closers, fn)
	}

	chRunDone := make(chan *component, total)
	signal.Notify(a.shutdown, sig...)

	var (
		shutdownRunning   bool
		notBlockedRunning int
	)

	defer func() {
		if notBlockedRunning >= total {
			close(chRunDone)
		}

		close(a.shutdown)
	}()

	// запускаем компоненты
	for _, cmp := range components {
		go func(component *component) {
			component.Run(a)
			chRunDone <- component
		}(cmp)
	}

	for {
		select {
		case cmp := <-chRunDone:
			notBlockedRunning++

			// если Run вернул ошибку, то завершаем все приложение
			if cmp.Status() == ComponentStatusRunFailed {
				if !shutdownRunning {
					closers = append([]func() error{
						func() error {
							return errors.New("component " + cmp.Name() + " run failed with error: " + cmp.RunError().Error())
						},
					}, closers...)

					a.shutdown <- sig[0]
				}
			}

			if notBlockedRunning >= total {
				if !shutdownRunning {
					a.shutdown <- sig[0]
				}
			}

		case <-a.shutdown:
			shutdownRunning = true // nolint:ineffassign
			signal.Stop(a.shutdown)

			var shutdownEG errgroup.Group

			for _, closer := range closers {
				shutdownEG.Go(closer)
			}

			return shutdownEG.Wait()
		}
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

func (a *App) GetComponent(n string) Component {
	if cmp, ok := a.components.Get(n); ok {
		return cmp.instance
	}

	return nil
}

func (a *App) GetComponents() ([]Component, error) {
	components, err := a.components.All()
	if err != nil {
		return nil, err
	}

	resolveComponents := make([]Component, len(components))
	for _, cmp := range components {
		resolveComponents[cmp.Order()] = cmp.instance
	}

	return resolveComponents, nil
}

func (a *App) HasComponent(n string) bool {
	return a.GetComponent(n) != nil
}

func (a *App) RegisterComponent(c Component) error {
	if a.HasComponent(c.Name()) {
		return errors.New("component \"" + c.Name() + "\" already exists")
	}

	a.components.Add(c.Name(), c)
	return nil
}

func (a *App) MustRegisterComponent(c Component) {
	if err := a.RegisterComponent(c); err != nil {
		panic(err)
	}
}

// nolint:golint
func (a *App) StatusComponent(n string) componentStatus {
	if cmp, ok := a.components.Get(n); ok {
		return cmp.Status()
	}

	return ComponentStatusUnknown
}

func (a *App) WatchComponentStatus(status componentStatus, name string, names ...string) <-chan struct{} {
	ch := make(chan struct{}, 1)
	ns := append([]string{name}, names...)

	if len(ns) == 0 {
		close(ch)
		return ch
	}

	go func() {
		for _, n := range ns {
			if cmp, ok := a.components.Get(n); ok {
				<-cmp.WatchStatus(status)
			}
		}

		ch <- struct{}{}
	}()

	return ch
}

func (a *App) ReadyComponent(name string, names ...string) <-chan struct{} {
	return a.WatchComponentStatus(ComponentStatusReady, name, names...)
}

func (a *App) RunningComponent(name string, names ...string) <-chan struct{} {
	return a.WatchComponentStatus(ComponentStatusFinished, name, names...)
}

func (a *App) ShutdownComponent(name string, names ...string) <-chan struct{} {
	return a.WatchComponentStatus(ComponentStatusShutdown, name, names...)
}

func (a *App) Shutdown() error {
	if atomic.LoadInt64(&a.running) != 1 {
		return errors.New("already shutdown")
	}

	a.shutdown <- sig[0]
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
