package frontend

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
)

type FrontendMenu struct {
	Name    string
	Url     string
	Icon    string
	SubMenu []*FrontendMenu
}

type ServiceFrontendHandlers interface {
	SetFrontendHandlers(*Router)
}

type ServiceFrontendMenu interface {
	GetFrontendMenu() *FrontendMenu
}

type FrontendService struct {
	Logger      *logrus.Entry
	config      *resource.Config
	template    *resource.Template
	application *shadow.Application
	router      *Router
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
	s.Logger = logger.Get(s.GetName())

	// скидывает mux по-умолчанию, так как pprof добавил свои хэндлеры
	http.DefaultServeMux = http.NewServeMux()

	s.router = NewRouter(s.application, s.Logger, s.config)
	s.router.SetPanicHandler(s, &PanicHandler{})
	s.router.SetNotAllowedHandler(s, &MethodNotAllowedHandler{})
	s.router.SetNotFoundHandler(s, &NotFoundHandler{})

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

	s.initAlerts()

	return nil
}

func (s *FrontendService) Run(wg *sync.WaitGroup) error {
	menus := make([]*FrontendMenu, 0, len(s.application.GetServices()))

	for _, service := range s.application.GetServices() {
		if serviceHandlers, ok := service.(ServiceFrontendHandlers); ok {
			serviceHandlers.SetFrontendHandlers(s.router)
		}

		if serviceMenu, ok := service.(ServiceFrontendMenu); ok {
			menu := serviceMenu.GetFrontendMenu()
			if menu != nil {
				if service == s {
					menus = append([]*FrontendMenu{menu}, menus...)
				} else {
					menus = append(menus, menu)
				}
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
		s.Logger.WithFields(fields).Info("Running service")

		if err := http.ListenAndServe(addr, s.router); err != nil {
			s.Logger.Fatalf("Could not start frontend [%d]: %s\n", os.Getpid(), err.Error())
		}
	}(s.router)

	return nil
}
