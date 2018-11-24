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
	GetComponent(string) Component
	GetComponents() ([]Component, error)
	HasComponent(string) bool
	RegisterComponent(Component) error
	IsReadyComponent(string) bool
	ReadyComponent(string) <-chan struct{}
}

type App struct {
	_ [4]byte // atomic requires 64-bit alignment for struct field access

	components *components
	running    int64

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
	}

	return application
}

func (a *App) Run() (err error) {
	if atomic.LoadInt64(&a.running) == 1 {
		return errors.New("already running")
	}

	atomic.StoreInt64(&a.running, 1)
	defer atomic.StoreInt64(&a.running, 0)

	components, err := a.components.all()
	if err != nil {
		return err
	}

	total := len(components)
	if total == 0 {
		return
	}

	closers := make([]func() error, 0, total)
	for _, cmp := range components {
		if closer, ok := cmp.instance.(ComponentShutdown); ok {
			closers = append(closers, closer.Shutdown)
		}
	}

	chResults := make(chan *component, total)

	chShutdown := make(chan os.Signal, 1)
	signal.Notify(chShutdown, sig...)

	shutdown := false

	defer func() {
		close(chResults)
		close(chShutdown)
	}()

	for _, cmp := range components {
		go cmp.Run(a, chResults)
	}

	var done int

	for {
		select {
		case cmp := <-chResults:
			if cmp.Done() {
				done++
			}

			if cmp.Error() != nil {
				if !shutdown {
					closers = append([]func() error{
						func() error {
							return cmp.Error()
						},
					}, closers...)

					chShutdown <- sig[0]
				}
			} else if done >= total {
				chShutdown <- sig[0]
			}

		case <-chShutdown:
			shutdown = true
			signal.Stop(chShutdown)

			var shutdownEG errgroup.Group

			for _, closer := range closers {
				shutdownEG.Go(closer)
			}

			return shutdownEG.Wait()
		}
	}

	return nil
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
	if cmp, ok := a.components.get(n); ok {
		return cmp.instance
	}

	return nil
}

func (a *App) GetComponents() ([]Component, error) {
	components, err := a.components.all()
	if err != nil {
		return nil, err
	}

	resolveComponents := make([]Component, len(components), len(components))
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

	if ini, ok := c.(ComponentInit); ok {
		if err := ini.Init(a); err != nil {
			return err
		}
	}

	a.components.add(c.Name(), c)
	return nil
}

func (a *App) MustRegisterComponent(c Component) {
	if err := a.RegisterComponent(c); err != nil {
		panic(err)
	}
}

func (a *App) IsReadyComponent(n string) bool {
	if cmp, ok := a.components.get(n); ok {
		return cmp.Ready()
	}

	return false
}

func (a *App) ReadyComponent(n string) <-chan struct{} {
	if cmp, ok := a.components.get(n); ok {
		return cmp.Watch()
	}

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
