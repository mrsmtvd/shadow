package frontend

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/justinas/alice"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
)

type FrontendMenu struct {
	Name    string
	Url     string
	SubMenu []*FrontendMenu
}

type ServiceFrontendHandlers interface {
	SetFrontendHandlers(*Router)
	GetFrontendMenu() *FrontendMenu
}

type FrontendService struct {
	config      *resource.Config
	template    *resource.Template
	application *shadow.Application
	router      *Router
	logger      *logrus.Entry
	middleware  alice.Chain
}

func (c *FrontendService) GetName() string {
	return "frontend"
}

func (s *FrontendService) Init(a *shadow.Application) (err error) {
	s.application = a

	resourceTemplate, err := a.GetResource("template")
	if err != nil {
		return err
	}
	s.template = resourceTemplate.(*resource.Template)
	s.template.Globals["Menu"] = make([]*FrontendMenu, 0, len(s.application.GetServices()))

	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	s.config = resourceConfig.(*resource.Config)

	resourceLogger, err := a.GetResource("logger")
	if err != nil {
		return err
	}
	logger := resourceLogger.(*resource.Logger)
	s.logger = logger.Get(s.GetName())

	// скидывает mux по-умолчанию, так как pprof добавил свои хэндлеры
	http.DefaultServeMux = http.NewServeMux()

	s.middleware = alice.New(
		LoggerMiddleware(s.logger),
		BasicAuthMiddleware(s.config.GetString("frontend.auth-user"), s.config.GetString("frontend.auth-password")),
	)
	s.router = NewRouter(s.application)

	panicHandler := &PanicHandler{}
	panicHandler.Init(s.application, s)
	s.router.PanicHandler = func(out http.ResponseWriter, in *http.Request, error interface{}) {
		panicHandler.InitRequest(out, in)
		panicHandler.SetError(error)
		panicHandler.Handle()
		panicHandler.Render()
	}

	methodNotAllowedHandler := MethodNotAllowedHandler{}
	methodNotAllowedHandler.Init(s.application, s)
	methodNotAllowed := func(out http.ResponseWriter, in *http.Request) {
		methodNotAllowedHandler.InitRequest(out, in)
		methodNotAllowedHandler.Handle()
		methodNotAllowedHandler.Render()
	}
	s.router.MethodNotAllowed = http.HandlerFunc(methodNotAllowed)

	notFoundHandler := &NotFoundHandler{}
	notFoundHandler.Init(s.application, s)
	notFound := func(out http.ResponseWriter, in *http.Request) {
		notFoundHandler.InitRequest(out, in)
		notFoundHandler.Handle()
		notFoundHandler.Render()
	}
	s.router.NotFound = http.HandlerFunc(notFound)

	if s.config.GetBool("debug") {
		s.router.HandlerFunc("GET", "/debug/pprof/cmdline", pprof.Cmdline)
		s.router.HandlerFunc("GET", "/debug/pprof/profile", pprof.Profile)
		s.router.HandlerFunc("GET", "/debug/pprof/symbol", pprof.Symbol)
		s.router.HandlerFunc("POST", "/debug/pprof/symbol", pprof.Symbol)
		s.router.HandlerFunc("GET", "/debug/pprof/block", pprof.Index)
		s.router.HandlerFunc("GET", "/debug/pprof/goroutine", pprof.Index)
		s.router.HandlerFunc("GET", "/debug/pprof/heap", pprof.Index)
		s.router.HandlerFunc("GET", "/debug/pprof/threadcreate", pprof.Index)
		s.router.HandlerFunc("GET", "/debug/pprof/", pprof.Index)
	}

	return nil
}

func (s *FrontendService) Run(wg *sync.WaitGroup) error {
	menus := make([]*FrontendMenu, 0, len(s.application.GetServices()))
	for _, service := range s.application.GetServices() {
		if serviceCast, ok := service.(ServiceFrontendHandlers); ok {
			serviceCast.SetFrontendHandlers(s.router)

			menu := serviceCast.GetFrontendMenu()
			if menu != nil {
				menus = append(menus, menu)
			}
		}
	}

	s.template.Globals["Menu"] = menus

	go func(router *Router) {
		defer wg.Done()

		http.HandleFunc("/", func(out http.ResponseWriter, in *http.Request) {
			router.ServeHTTP(out, in)
		})

		// TODO: ssl

		addr := fmt.Sprintf("%s:%d", s.config.GetString("frontend.host"), s.config.GetInt64("frontend.port"))
		fields := logrus.Fields{
			"addr": addr,
			"pid":  os.Getpid(),
		}
		s.logger.WithFields(fields).Info("Running service")

		if err := http.ListenAndServe(addr, s.middleware.Then(http.DefaultServeMux)); err != nil {
			s.logger.Fatalf("Could not start frontend [%d]: %s\n", os.Getpid(), err.Error())
		}
	}(s.router)

	return nil
}
