package frontend

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
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
	boxStatic   *rice.Box
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

	rice.Debug = s.config.GetBool("debug")
	s.boxStatic, err = rice.FindBox("public")
	if err != nil {
		return errors.Wrap(err, "Failed init frontend service")
	}

	// скидывает mux по-умолчанию, так как pprof добавил свои хэндлеры
	http.DefaultServeMux = http.NewServeMux()

	s.middleware = alice.New(
		LoggerMiddleware(s.logger),
		BasicAuthMiddleware(s.config.GetString("auth-user"), s.config.GetString("auth-password")),
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
	s.router.MethodNotAllowed = func(out http.ResponseWriter, in *http.Request) {
		methodNotAllowedHandler.InitRequest(out, in)
		methodNotAllowedHandler.Handle()
		methodNotAllowedHandler.Render()
	}

	notFoundHandler := &NotFoundHandler{}
	notFoundHandler.Init(s.application, s)
	s.router.NotFound = func(out http.ResponseWriter, in *http.Request) {
		notFoundHandler.InitRequest(out, in)
		notFoundHandler.Handle()
		notFoundHandler.Render()
	}

	if rice.Debug {
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

func (s *FrontendService) Run() error {
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
		http.HandleFunc("/", func(out http.ResponseWriter, in *http.Request) {
			router.ServeHTTP(out, in)
		})

		addr := fmt.Sprintf("%s:%d", s.config.GetString("host"), s.config.GetInt64("port"))
		s.logger.Infof("running frontend [%d]: %s", os.Getpid(), addr)

		if err := http.ListenAndServe(addr, s.middleware.Then(http.DefaultServeMux)); err != nil {
			s.logger.Fatalf("could not start frontend [%d]: %s\n", os.Getpid(), err.Error())
		}
	}(s.router)

	return nil
}
