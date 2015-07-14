package api

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
	"gopkg.in/jcelliott/turnpike.v2"
)

type ServiceApiHandler interface {
	GetMethods() map[string]turnpike.MethodHandler
}

type ApiService struct {
	application *shadow.Application
	config      *resource.Config
	logger      *logrus.Entry
}

func (s *ApiService) GetName() string {
	return "api"
}

func (s *ApiService) Init(a *shadow.Application) error {
	s.application = a

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

	return nil
}

func (s *ApiService) Run(wg *sync.WaitGroup) error {
	if s.config.GetBool("debug") {
		turnpike.Debug()
	}

	handler := turnpike.NewBasicWebsocketServer(s.GetName())
	client, err := handler.GetLocalClient(s.GetName())
	if err != nil {
		return err
	}

	for _, service := range s.application.GetServices() {
		if serviceCast, ok := service.(ServiceApiHandler); ok {
			for procedure, fn := range serviceCast.GetMethods() {
				procedure = fmt.Sprintf("%s.%s", service.GetName(), procedure)
				if err = client.Register(procedure, fn); err != nil {
					return err
				}
			}
		}
	}

	go func(handler *turnpike.WebsocketServer) {
		defer wg.Done()

		// TODO: ssl

		addr := fmt.Sprintf("%s:%d", s.config.GetString("api-host"), s.config.GetInt64("api-port"))
		fields := logrus.Fields{
			"addr": addr,
			"pid":  os.Getpid(),
		}
		s.logger.WithFields(fields).Info("Running service")

		mux := http.NewServeMux()
		server := &http.Server{
			Handler: mux,
			Addr:    addr,
		}

		mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		mux.Handle("/", handler)

		if err := server.ListenAndServe(); err != nil {
			s.logger.Fatalf("Could not start api [%d]: %s\n", os.Getpid(), err.Error())
		}
	}(handler)

	return nil
}
