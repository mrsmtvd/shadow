package grpc

import (
	"reflect"

	ptypes_struct "github.com/golang/protobuf/ptypes/struct"
)

func ConvertMapStringInterfaceToStructProto(m map[string]interface{}) *ptypes_struct.Struct {
	s := &ptypes_struct.Struct{
		Fields: make(map[string]*ptypes_struct.Value, len(m)),
	}

	for k, v := range m {
		s.Fields[k] = ConvertInterfaceToStructValueProto(v)
	}

	return s
}

func ConvertStructProtoToMapStringInterface(s *ptypes_struct.Struct) map[string]interface{} {
	m := make(map[string]interface{}, len(s.GetFields()))

	for key, field := range s.GetFields() {
		m[key] = ConvertStructValueProtoToInterface(field)
	}

	return m
}

func ConvertStructValueProtoToInterface(value *ptypes_struct.Value) interface{} {
	switch cast := value.GetKind().(type) {
	case *ptypes_struct.Value_NullValue:
		return nil

	case *ptypes_struct.Value_NumberValue:
		return cast.NumberValue

	case *ptypes_struct.Value_StringValue:
		return cast.StringValue

	case *ptypes_struct.Value_BoolValue:
		return cast.BoolValue

	case *ptypes_struct.Value_StructValue:
		list := make(map[string]interface{}, len(cast.StructValue.GetFields()))

		for key, listValue := range cast.StructValue.GetFields() {
			list[key] = ConvertStructValueProtoToInterface(listValue)
		}

		return list

	case *ptypes_struct.Value_ListValue:
		list := make([]interface{}, 0, len(cast.ListValue.GetValues()))

		for _, listValue := range cast.ListValue.GetValues() {
			list = append(list, ConvertStructValueProtoToInterface(listValue))
		}

		return list
	}

	return nil
}

func ConvertInterfaceToStructValueProto(v interface{}) *ptypes_struct.Value {
	switch cast := v.(type) {
	case nil:
		break

	case bool:
		return &ptypes_struct.Value{
			Kind: &ptypes_struct.Value_BoolValue{
				BoolValue: cast,
			},
		}

	case []byte:
		return ConvertInterfaceToStructValueProto(string(cast))

	case string:
		return &ptypes_struct.Value{
			Kind: &ptypes_struct.Value_StringValue{
				StringValue: cast,
			},
		}

	case int:
		return ConvertInterfaceToStructValueProto(float64(cast))

	case int8:
		return ConvertInterfaceToStructValueProto(float64(cast))

	case int16:
		return ConvertInterfaceToStructValueProto(float64(cast))

	case int32:
		return ConvertInterfaceToStructValueProto(float64(cast))

	case int64:
		return ConvertInterfaceToStructValueProto(float64(cast))

	case uint:
		return ConvertInterfaceToStructValueProto(float64(cast))

	case uint8:
		return ConvertInterfaceToStructValueProto(float64(cast))

	case uint16:
		return ConvertInterfaceToStructValueProto(float64(cast))

	case uint32:
		return ConvertInterfaceToStructValueProto(float64(cast))

	case uint64:
		return ConvertInterfaceToStructValueProto(float64(cast))

	case float32:
		return ConvertInterfaceToStructValueProto(float64(cast))

	case float64:
		return &ptypes_struct.Value{
			Kind: &ptypes_struct.Value_NumberValue{
				NumberValue: cast,
			},
		}

		// case complex64:

		// case complex128:

	default:
		value := reflect.ValueOf(v)

		switch value.Kind() {
		case reflect.Array, reflect.Slice:
			items := make([]*ptypes_struct.Value, 0, value.Len())

			for i := 0; i < value.Len(); i++ {
				items = append(items, ConvertInterfaceToStructValueProto(value.Index(i).Interface()))
			}

			return &ptypes_struct.Value{
				Kind: &ptypes_struct.Value_ListValue{
					ListValue: &ptypes_struct.ListValue{
						Values: items,
					},
				},
			}

		case reflect.Map:
			items := make(map[string]*ptypes_struct.Value, value.Len())

			for _, n := range value.MapKeys() {
				items[n.String()] = ConvertInterfaceToStructValueProto(value.MapIndex(n).Interface())
			}

			return &ptypes_struct.Value{
				Kind: &ptypes_struct.Value_StructValue{
					StructValue: &ptypes_struct.Struct{
						Fields: items,
					},
				},
			}

		case reflect.Struct:
			items := make(map[string]*ptypes_struct.Value, value.NumField())

			for i := 0; i < value.NumField(); i++ {
				field := value.Type().Field(i)
				items[field.Name] = ConvertInterfaceToStructValueProto(value.FieldByName(field.Name).Interface())
			}

			return &ptypes_struct.Value{
				Kind: &ptypes_struct.Value_StructValue{
					StructValue: &ptypes_struct.Struct{
						Fields: items,
					},
				},
			}

		case reflect.Ptr:
			return ConvertInterfaceToStructValueProto(reflect.Indirect(reflect.ValueOf(v)).Interface())

			//case reflect.Uintptr, reflect.UnsafePointer, reflect.Chan, reflect.Func, reflect.Interface, reflect.Invalid:
			//	return ConvertInterfaceToStructValueProto(value.String())
		}
	}

	return nil
}
