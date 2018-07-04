package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/grpc"
	"golang.org/x/net/context"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

// easyjson:json
type managerHandlerResponseCall struct {
	Result   string      `json:"result,omitempty"`
	Error    string      `json:"error,omitempty"`
	Headers  metadata.MD `json:"headers,omitempty"`
	Trailers metadata.MD `json:"trailers,omitempty"`
}

type managerHandlerServiceViewData struct {
	Name    string
	Methods []*ManagerHandlerMethodViewData
}

type ManagerHandlerMethodViewData struct {
	Name         string
	InputStream  bool
	OutputStream bool
	InputType    *ManagerHandlerMessageViewData
	OutputType   *ManagerHandlerMessageViewData
}

type ManagerHandlerMessageViewData struct {
	Name   string
	Fields []*ManagerHandlerFieldViewData
}

type ManagerHandlerFieldViewData struct {
	Name        string
	Type        string
	Default     interface{}
	Label       string
	Enum        []*ManagerHandlerFieldViewDataEnum
	Message     *ManagerHandlerMessageViewData
	IsEnum      bool
	IsExtension bool
	IsMap       bool
	IsMessage   bool
	IsRequired  bool
	IsRepeated  bool
}

type ManagerHandlerFieldViewDataEnum struct {
	Number    int32
	Name      string
	IsDefault bool
}

type ManagerHandler struct {
	dashboard.Handler
}

func getTypeName(field *desc.FieldDescriptor) string {
	var name string

	if field.IsMap() {
		name = fmt.Sprintf("map<%s>%s", getTypeName(field.GetMapKeyType()), getTypeName(field.GetMapValueType()))
	} else if field.GetMessageType() != nil {
		name = field.GetMessageType().GetName()
	} else if field.GetEnumType() != nil {
		name = field.GetEnumType().GetName()
	} else {
		name = strings.ToLower(field.GetType().String()[5:])
	}

	if field.IsRepeated() {
		name = "[]" + name
	}

	return name
}

func getMessageViewDate(message *desc.MessageDescriptor, currentLevel, maxLevel int64) *ManagerHandlerMessageViewData {
	fields := message.GetFields()

	view := &ManagerHandlerMessageViewData{
		Name:   message.GetName(),
		Fields: make([]*ManagerHandlerFieldViewData, len(fields), len(fields)),
	}

	for i, f := range fields {
		view.Fields[i] = getFieldViewDate(f, currentLevel, maxLevel)
	}

	return view
}

func getFieldViewDate(field *desc.FieldDescriptor, currentLevel, maxLevel int64) *ManagerHandlerFieldViewData {
	data := &ManagerHandlerFieldViewData{
		Name:        field.GetName(),
		Type:        getTypeName(field),
		Default:     field.GetDefaultValue(),
		Label:       strings.ToLower(field.GetLabel().String()[6:]),
		IsExtension: field.IsExtension(),
		IsEnum:      field.GetEnumType() != nil,
		IsMap:       field.IsMap(),
		IsMessage:   field.GetMessageType() != nil,
		IsRequired:  field.IsRequired(),
		IsRepeated:  field.IsRepeated(),
	}

	if field.GetMessageType() != nil && currentLevel <= maxLevel {
		data.Message = getMessageViewDate(field.GetMessageType(), currentLevel+1, maxLevel)
	}

	// Scalar Value Types
	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SINT64,
		descriptor.FieldDescriptorProto_TYPE_FLOAT,
		descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		data.Default = 0
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		data.Default = false
	case descriptor.FieldDescriptorProto_TYPE_BYTES,
		descriptor.FieldDescriptorProto_TYPE_STRING:
		data.Default = ""
	}

	// Enumerations
	if field.GetEnumType() != nil {
		values := field.GetEnumType().GetValues()
		data.Enum = make([]*ManagerHandlerFieldViewDataEnum, len(values), len(values))

		for e, enum := range values {
			def, ok := field.GetDefaultValue().(int32)

			data.Enum[e] = &ManagerHandlerFieldViewDataEnum{
				Number:    enum.GetNumber(),
				Name:      enum.GetName(),
				IsDefault: ok && def == enum.GetNumber(),
			}
		}
	}

	return data
}

func (v *ManagerHandlerFieldViewData) MarshalJSON() ([]byte, error) {
	var d interface{}

	if v.Message != nil {
		s := make(map[string]interface{}, len(v.Message.Fields))

		for _, f := range v.Message.Fields {
			s[f.Name] = f
		}

		d = s
	} else if v.IsMap {
		d = map[string]interface{}(nil)
	} else {
		d = v.Default
	}

	/*
		if v.IsRepeated {
			d = []interface{}{d}
		}
	*/

	return json.Marshal(d)
}

func (v *ManagerHandlerFieldViewData) JSON() string {
	if j, err := json.MarshalIndent(v, "", "    "); err == nil {
		return string(j)
	}

	return ""
}

func (h *ManagerHandler) getServicesLightViewData(cmp grpc.Component) ([]managerHandlerServiceViewData, error) {
	list := cmp.GetServiceInfo()
	ret := make([]managerHandlerServiceViewData, 0, len(list))

	for name, info := range list {
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

func (h *ManagerHandler) getServicesViewData(r *dashboard.Request) ([]managerHandlerServiceViewData, error) {
	ret := make([]managerHandlerServiceViewData, 0)

	addr := net.JoinHostPort(r.Config().String(grpc.ConfigHost), r.Config().String(grpc.ConfigPort))
	connect, err := g.DialContext(r.Context(), addr, g.WithInsecure())
	if err != nil {
		return ret, err
	}

	cli := grpcreflect.NewClient(r.Context(), rpb.NewServerReflectionClient(connect))
	maxLevel := r.Config().Int64(grpc.ConfigManagerMaxLevel)

	if services, err := cli.ListServices(); err == nil {
		for _, s := range services {
			service, err := cli.ResolveService(s)
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
					InputType:    getMessageViewDate(m.GetInputType(), 1, maxLevel),
					OutputType:   getMessageViewDate(m.GetOutputType(), 1, maxLevel),
				})
			}

			ret = append(ret, view)
		}
	}

	return ret, nil
}

func (h *ManagerHandler) actionCall(w *dashboard.Response, r *dashboard.Request) {
	s := r.Original().FormValue("_service")
	m := r.Original().FormValue("_method")

	if s == "" || m == "" {
		h.NotFound(w, r)
		return
	}

	addr := net.JoinHostPort(r.Config().String(grpc.ConfigHost), r.Config().String(grpc.ConfigPort))
	connect, err := g.DialContext(r.Context(), addr, g.WithInsecure())
	if err != nil {
		w.SendJSON(managerHandlerResponseCall{
			Error: err.Error(),
		})
		return
	}

	cli := grpcreflect.NewClient(r.Context(), rpb.NewServerReflectionClient(connect))

	service, err := cli.ResolveService(s)
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

	stub := grpcdynamic.NewStub(connect)
	ctx := context.Background()
	request := dynamic.NewMessage(method.GetInputType())

	var (
		headerMD  metadata.MD
		trailerMD metadata.MD
	)

	result, err := stub.InvokeRpc(ctx, method, request, g.Header(&headerMD), g.Trailer(&trailerMD))
	response := managerHandlerResponseCall{
		Headers:  headerMD,
		Trailers: trailerMD,
	}

	if err != nil {
		response.Error = err.Error()

		st, ok := status.FromError(err)
		if ok {
			result = st.Proto()
		}
	}

	marshaler := &jsonpb.Marshaler{}

	if result != nil {
		resultJSON, err := marshaler.MarshalToString(result)
		if err == nil {
			response.Result = resultJSON
		}
	}

	w.SendJSON(response)
}

func (h *ManagerHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if r.IsPost() {
		if r.URL().Query().Get("action") == "call" {
			h.actionCall(w, r)
		} else {
			h.NotFound(w, r)
		}

		return
	}

	var (
		err      error
		services []managerHandlerServiceViewData
	)

	if r.Config().Bool(grpc.ConfigReflectionEnabled) {
		services, err = h.getServicesViewData(r)
	} else {
		services, err = h.getServicesLightViewData(r.Component().(grpc.Component))
	}

	h.Render(r.Context(), "manager", map[string]interface{}{
		"error":    err,
		"services": services,
	})
}
