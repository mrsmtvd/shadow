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
type indexHandlerResponseCall struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

type IndexHandlerServiceViewData struct {
	Name    string
	Methods []*IndexHandlerMethodViewData
}

type IndexHandlerMethodViewData struct {
	Name         string
	InputStream  bool
	OutputStream bool
	InputType    *IndexHandlerTypeViewData
	OutputType   *IndexHandlerTypeViewData
}

type IndexHandlerTypeViewData struct {
	Name   string
	Fields []*IndexHandlerFieldViewData
}

type IndexHandlerFieldViewData struct {
	Name   string
	Type   string
	Labels []string
}

type IndexHandler struct {
	dashboard.Handler

	connect *g.ClientConn
	cli     *grpcreflect.Client
	server  *g.Server
}

func NewIndexHandler(c config.Component, s *g.Server) *IndexHandler {
	h := &IndexHandler{
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

func (h *IndexHandler) getServicesLightViewData() ([]IndexHandlerServiceViewData, error) {
	ret := []IndexHandlerServiceViewData{}

	for name, info := range h.server.GetServiceInfo() {
		if proto.FileDescriptor(info.Metadata.(string)) == nil {
			continue
		}

		view := IndexHandlerServiceViewData{
			Name:    name,
			Methods: []*IndexHandlerMethodViewData{},
		}

		for _, method := range info.Methods {
			view.Methods = append(view.Methods, &IndexHandlerMethodViewData{
				Name: method.Name,
			})
		}

		ret = append(ret, view)
	}

	return ret, nil
}

func (h *IndexHandler) getServicesViewData() ([]IndexHandlerServiceViewData, error) {
	fillType := func(t *desc.MessageDescriptor) *IndexHandlerTypeViewData {
		fields := t.GetFields()

		view := &IndexHandlerTypeViewData{
			Name:   t.GetName(),
			Fields: make([]*IndexHandlerFieldViewData, len(fields), len(fields)),
		}

		for i, f := range fields {
			view.Fields[i] = &IndexHandlerFieldViewData{
				Name: f.GetName(),
				Type: strings.ToLower(f.GetType().String()[5:]),
			}
		}

		return view
	}

	ret := []IndexHandlerServiceViewData{}
	if services, err := h.cli.ListServices(); err == nil {
		for _, s := range services {
			service, err := h.cli.ResolveService(s)
			if err != nil {
				return ret, err
			}

			view := IndexHandlerServiceViewData{
				Name:    s,
				Methods: []*IndexHandlerMethodViewData{},
			}

			for _, m := range service.GetMethods() {
				view.Methods = append(view.Methods, &IndexHandlerMethodViewData{
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

func (h *IndexHandler) call(w *dashboard.Response, r *dashboard.Request) {
	s := r.Original().FormValue("service")
	m := r.Original().FormValue("method")

	if s == "" || m == "" {
		h.NotFound(w, r)
		return
	}

	service, err := h.cli.ResolveService(s)
	if err != nil {
		w.SendJSON(indexHandlerResponseCall{
			Error: err.Error(),
		})
		return
	}

	method := service.FindMethodByName(m)
	if method == nil {
		w.SendJSON(indexHandlerResponseCall{
			Error: "Method not found",
		})
		return
	}

	stub := grpcdynamic.NewStub(h.connect)
	ctx := context.Background()
	request := dynamic.NewMessage(method.GetInputType())

	result, err := stub.InvokeRpc(ctx, method, request)
	if err != nil {
		w.SendJSON(indexHandlerResponseCall{
			Error: err.Error(),
		})
		return
	}

	marshaler := &jsonpb.Marshaler{}
	responseJSON, err := marshaler.MarshalToString(result)
	if err != nil {
		w.SendJSON(indexHandlerResponseCall{
			Error: err.Error(),
		})
		return
	}

	w.SendJSON(indexHandlerResponseCall{
		Result: responseJSON,
	})
}

func (h *IndexHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
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
		services []IndexHandlerServiceViewData
	)

	if r.Config().GetBool(grpc.ConfigReflectionEnabled) {
		services, err = h.getServicesViewData()
	} else {
		services, err = h.getServicesLightViewData()
	}

	h.Render(r.Context(), grpc.ComponentName, "index", map[string]interface{}{
		"error":    err,
		"services": services,
	})
}
