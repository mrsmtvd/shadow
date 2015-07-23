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
	GetApiProcedures() []ApiProcedure
}

type ApiService struct {
	application *shadow.Application
	config      *resource.Config
	logger      *logrus.Entry
	procedures  []ApiProcedure
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
			for _, procedure := range serviceCast.GetApiProcedures() {
				name := procedure.GetName()

				logEntry := s.logger.WithFields(logrus.Fields{
					"procedure": name,
					"service":   service.GetName(),
				})

				if s.HasProcedure(name) {
					logEntry.Warn("Procedure already exists. Ignore procedure.")
					continue
				}

				procedure.Init(s, s.application)
				procedureWrapper := func(procedure ApiProcedure) turnpike.MethodHandler {
					return func(args []interface{}, kwargs map[string]interface{}) *turnpike.CallResult {
						if validator, ok := procedure.(ApiProcedureRequest); ok {
							request := validator.GetRequest()
							if err := RequestFillAndValidate(request, args, kwargs); err != nil {
								return &turnpike.CallResult{
									Err: turnpike.URI(err.Error()),
								}
							}

							return validator.Run(request)
						}

						if simple, ok := procedure.(ApiProcedureSimple); ok {
							return simple.Run(args, kwargs)
						}

						logEntry.WithField("error", err.Error()).Error("Error procedure interace")
						return &turnpike.CallResult{
							Err: turnpike.URI(ErrorUnknownProcedure),
						}
					}
				}

				if err = client.Register(name, procedureWrapper(procedure)); err != nil {
					logEntry.WithField("error", err.Error()).Error("Error register api procedure")
					// ignore error
				} else {
					logEntry.Debug("Register procedure")
				}
				s.procedures = append(s.procedures, procedure)
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
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// FiXME: Magic
			delete(r.Header, "Origin")

			s.logger.Infof("Connection from %s", r.RemoteAddr)
			handler.ServeHTTP(w, r)
		})

		if err := server.ListenAndServe(); err != nil {
			s.logger.Fatalf("Could not start api [%d]: %s\n", os.Getpid(), err.Error())
		}
	}(handler)

	return nil
}

func (s *ApiService) GetProcedures() []ApiProcedure {
	return s.procedures
}

func (s *ApiService) HasProcedure(procedure string) bool {
	for _, p := range s.procedures {
		if p.GetName() == procedure {
			return true
		}
	}

	return false
}

func (s *ApiService) GetClient() (*turnpike.Client, error) {
	addr := fmt.Sprintf("ws://%s:%d/", s.config.GetString("api-host"), s.config.GetInt64("api-port"))

	client, err := turnpike.NewWebsocketClient(turnpike.JSON, addr)
	if err != nil {
		return nil, err
	}

	_, err = client.JoinRealm("api", turnpike.ALLROLES, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
