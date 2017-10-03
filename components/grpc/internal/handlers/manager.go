package handlers

import (
	"net"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/grpc"
	"golang.org/x/net/context"
	g "google.golang.org/grpc"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// easyjson:json
type managerHandlerResponseCall struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

type managerHandlerServiceViewData struct {
	Name    string
	Methods []*ManagerHandlerMethodViewData
}

type ManagerHandlerMethodViewData struct {
	Name         string
	InputStream  bool
	OutputStream bool
	InputType    *ManagerHandlerTypeViewData
	OutputType   *ManagerHandlerTypeViewData
}

type ManagerHandlerTypeViewData struct {
	Name   string
	Fields []*ManagerHandlerFieldViewData
}

type ManagerHandlerFieldViewData struct {
	Name   string
	Type   string
	Labels []string
}

type ManagerHandler struct {
	dashboard.Handler

	connect *g.ClientConn
	cli     *grpcreflect.Client
	server  *g.Server
}

func NewManagerHandler(c config.Component, s *g.Server) *ManagerHandler {
	h := &ManagerHandler{
		server: s,
	}

	ctx := context.Background()
	addr := net.JoinHostPort(c.GetString(grpc.ConfigHost), c.GetString(grpc.ConfigPort))

	var err error

	if h.connect, err = g.DialContext(ctx, addr, g.WithInsecure()); err == nil {
		h.cli = grpcreflect.NewClient(ctx, rpb.NewServerReflectionClient(h.connect))
	}

	return h
}

func (h *ManagerHandler) getServicesLightViewData() ([]managerHandlerServiceViewData, error) {
	ret := []managerHandlerServiceViewData{}

	for name, info := range h.server.GetServiceInfo() {
		if proto.FileDescriptor(info.Metadata.(string)) == nil {
			continue
		}

		view := managerHandlerServiceViewData{
			Name:    name,
			Methods: []*ManagerHandlerMethodViewData{},
		}

		for _, method := range info.Methods {
			view.Methods = append(view.Methods, &ManagerHandlerMethodViewData{
				Name: method.Name,
			})
		}

		ret = append(ret, view)
	}

	return ret, nil
}

func (h *ManagerHandler) getServicesViewData() ([]managerHandlerServiceViewData, error) {
	fillType := func(t *desc.MessageDescriptor) *ManagerHandlerTypeViewData {
		fields := t.GetFields()

		view := &ManagerHandlerTypeViewData{
			Name:   t.GetName(),
			Fields: make([]*ManagerHandlerFieldViewData, len(fields), len(fields)),
		}

		for i, f := range fields {
			view.Fields[i] = &ManagerHandlerFieldViewData{
				Name: f.GetName(),
				Type: strings.ToLower(f.GetType().String()[5:]),
			}
		}

		return view
	}

	ret := []managerHandlerServiceViewData{}
	if services, err := h.cli.ListServices(); err == nil {
		for _, s := range services {
			service, err := h.cli.ResolveService(s)
			if err != nil {
				return ret, err
			}

			view := managerHandlerServiceViewData{
				Name:    s,
				Methods: []*ManagerHandlerMethodViewData{},
			}

			for _, m := range service.GetMethods() {
				view.Methods = append(view.Methods, &ManagerHandlerMethodViewData{
					Name:         m.GetName(),
					InputStream:  m.IsClientStreaming(),
					OutputStream: m.IsServerStreaming(),
					InputType:    fillType(m.GetInputType()),
					OutputType:   fillType(m.GetOutputType()),
				})
			}

			ret = append(ret, view)
		}
	}

	return ret, nil
}

func (h *ManagerHandler) call(w *dashboard.Response, r *dashboard.Request) {
	s := r.Original().FormValue("service")
	m := r.Original().FormValue("method")

	if s == "" || m == "" {
		h.NotFound(w, r)
		return
	}

	service, err := h.cli.ResolveService(s)
	if err != nil {
		w.SendJSON(managerHandlerResponseCall{
			Error: err.Error(),
		})
		return
	}

	method := service.FindMethodByName(m)
	if method == nil {
		w.SendJSON(managerHandlerResponseCall{
			Error: "Method not found",
		})
		return
	}

	stub := grpcdynamic.NewStub(h.connect)
	ctx := context.Background()
	request := dynamic.NewMessage(method.GetInputType())

	result, err := stub.InvokeRpc(ctx, method, request)
	if err != nil {
		w.SendJSON(managerHandlerResponseCall{
			Error: err.Error(),
		})
		return
	}

	marshaler := &jsonpb.Marshaler{}
	responseJSON, err := marshaler.MarshalToString(result)
	if err != nil {
		w.SendJSON(managerHandlerResponseCall{
			Error: err.Error(),
		})
		return
	}

	w.SendJSON(managerHandlerResponseCall{
		Result: responseJSON,
	})
}

func (h *ManagerHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if r.IsPost() {
		if r.URL().Query().Get("action") == "call" {
			h.call(w, r)
		} else {
			h.NotFound(w, r)
		}

		return
	}

	var (
		err      error
		services []managerHandlerServiceViewData
	)

	if r.Config().GetBool(grpc.ConfigReflectionEnabled) {
		services, err = h.getServicesViewData()
	} else {
		services, err = h.getServicesLightViewData()
	}

	h.Render(r.Context(), grpc.ComponentName, "manager", map[string]interface{}{
		"error":    err,
		"services": services,
	})
}
