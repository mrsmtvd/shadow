package api

import (
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/dropbox/godropbox/errors"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
)

type ServiceApiHandler interface {
	GetProcessor() thrift.TProcessor
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
	serverTransport, err := s.GetServerTransport()
	if err != nil {
		return err
	}
	defer serverTransport.Close()

	// protocol
	protocol := s.config.GetString("api-protocol")
	protocolFactory, err := s.GetProtocolFactory(s.config.GetString("api-protocol"))
	if err != nil {
		return err
	}

	if _, ok := protocolFactory.(*thrift.TDebugProtocolFactory); ok {
		protocol = "debug"
	}

	// transport
	transport := s.config.GetString("api-transport")
	transportFactory, err := s.GetTransportFactory(transport)
	if err != nil {
		return err
	}

	processor := thrift.NewTMultiplexedProcessor()
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	for _, service := range s.application.GetServices() {
		if serviceCast, ok := service.(ServiceApiHandler); ok {
			processor.RegisterProcessor(service.GetName(), serviceCast.GetProcessor())
		}
	}

	go server.Serve()

	fields := logrus.Fields{
		"protocol":  protocol,
		"transport": transport,
	}

	if socket, ok := serverTransport.(*thrift.TServerSocket); ok {
		fields["addr"] = socket.Addr()
	}

	s.logger.WithFields(fields).Info("Running service")

	return nil
}

func (s *ApiService) GetClientTransport() (thrift.TTransport, error) {
	addr := fmt.Sprintf("%s:%s", s.config.GetString("api-host"), s.config.GetString("api-port"))
	return thrift.NewTSocket(addr)
}

func (s *ApiService) GetServerTransport() (thrift.TServerTransport, error) {
	addr := fmt.Sprintf("%s:%s", s.config.GetString("api-host"), s.config.GetString("api-port"))
	return thrift.NewTServerSocket(addr)
}

func (s *ApiService) GetProtocolFactory(protocol string) (protocolFactory thrift.TProtocolFactory, err error) {
	switch protocol {
	case "binary":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	default:
		return nil, errors.Newf("Invalid protocol specified %s", protocol)
	}

	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	if s.config.GetBool("debug") {
		protocolFactory = thrift.NewTDebugProtocolFactory(protocolFactory, "shadow:")
	}

	return protocolFactory, nil
}

func (s *ApiService) GetTransportFactory(transport string) (thrift.TTransportFactory, error) {
	switch transport {
	case "buffered":
		return thrift.NewTBufferedTransportFactory(8192), nil
	case "framed":
		transportFactory := thrift.NewTTransportFactory()
		return thrift.NewTFramedTransportFactory(transportFactory), nil
	case "":
		return thrift.NewTTransportFactory(), nil
	default:
		return nil, errors.Newf("Invalid transport specified %s", transport)
	}
}
