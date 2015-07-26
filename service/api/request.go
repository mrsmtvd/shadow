package api

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kihamo/shadow"
)

const (
	ParamNameTag = "param"
)

type Request struct {
	zeros  map[string]bool
	errors []string
	out    interface{}
	args   []interface{}
	kwargs map[string]interface{}
}

func NewRequest(out interface{}, args []interface{}, kwargs map[string]interface{}) *Request {
	request := &Request{
		zeros:  map[string]bool{},
		errors: []string{},
		out:    out,
		args:   args,
		kwargs: kwargs,
	}
	outValue := reflect.ValueOf(out).Elem()

	if outValue.Kind() == reflect.Ptr && outValue.IsValid() {
		outValue = outValue.Elem()
	}

	switch outValue.Kind() {
	case reflect.Struct, reflect.Map, reflect.Interface:
		request.fill(outValue, kwargs, "")

	default:
		request.fill(outValue, args, "")
	}

	return request
}

func (r *Request) Valid() (bool, []string) {
	r.validateAny(reflect.ValueOf(r.out).Elem(), "")
	return len(r.errors) == 0, r.errors
}

func (r *Request) GetArgs() []interface{} {
	return r.args
}

func (r *Request) GetKwargs() map[string]interface{} {
	return r.kwargs
}

func (r *Request) isZero(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Slice, reflect.Map:
		return val.IsNil()

	case reflect.Struct:
		z := true
		for i := 0; i < val.NumField(); i++ {
			z = z && r.isZero(val.Field(i))
		}
		return z
	}

	return reflect.Zero(val.Type()).Interface() == val.Interface()
}

func (r *Request) setZero(val reflect.Value, path string) {
	if path != "" && r.isZero(val) {
		r.zeros[path] = true
	}
}

func (r *Request) getPath(parentPath string, childPath string) string {
	prefix := "."
	if parentPath == "" {
		prefix = ""
	}

	return parentPath + prefix + childPath
}

func (r *Request) fill(out reflect.Value, in interface{}, path string) {
	switch out.Kind() {
	case reflect.Ptr:
		if out.IsValid() {
			if out.IsNil() {
				out.Set(reflect.New(out.Type().Elem()))
			}

			r.fill(out.Elem(), in, path)
		}

	case reflect.Interface:
		r.fill(out.Elem(), in, path)

	case reflect.Struct:
		if in, ok := in.(map[string]interface{}); ok {
			r.setZero(out, path)

			for i := 0; i < out.NumField(); i++ {
				f := out.Type().Field(i)

				name := f.Tag.Get(ParamNameTag)
				if name == "" {
					name = strings.ToLower(f.Name)
				}

				childPath := r.getPath(path, name)
				if value, ok := in[name]; ok {
					r.fill(out.Field(i), value, childPath)
				}
			}
		}

	case reflect.Map:
		if in, ok := in.(map[string]interface{}); ok {
			r.setZero(out, path)

			keyType := out.Type().Key()
			valueType := out.Type().Elem()

			out.Set(reflect.MakeMap(out.Type()))
			for mapKey, mapValue := range in {
				key := reflect.New(keyType).Elem()
				r.fill(key, mapKey, path)

				value := reflect.New(valueType).Elem()
				childPath := r.getPath(path, fmt.Sprintf("[%q]", mapKey))
				r.fill(value, mapValue, childPath)

				out.SetMapIndex(key, value)
			}
		}

	case reflect.Slice:
		if in, ok := in.([]interface{}); ok {
			r.setZero(out, path)

			out.Set(reflect.MakeSlice(out.Type(), len(in), cap(in)))
			for i := range in {
				r.fill(out.Index(i), in[i], r.getPath(path, fmt.Sprintf("[%d]", i)))
			}
		}

	case reflect.Bool:
		r.setZero(out, path)
		out.SetBool(shadow.ToBool(in))

	case reflect.String:
		r.setZero(out, path)
		out.SetString(shadow.ToString(in))

	case reflect.Uint:
		r.setZero(out, path)
		out.Set(reflect.ValueOf(shadow.ToUint(in)))

	case reflect.Uint8:
		r.setZero(out, path)
		out.Set(reflect.ValueOf(shadow.ToUint8(in)))

	case reflect.Uint16:
		r.setZero(out, path)
		out.Set(reflect.ValueOf(shadow.ToUint16(in)))

	case reflect.Uint32:
		r.setZero(out, path)
		out.Set(reflect.ValueOf(shadow.ToUint32(in)))

	case reflect.Uint64:
		r.setZero(out, path)
		out.SetUint(shadow.ToUint64(in))

	case reflect.Int:
		r.setZero(out, path)
		out.Set(reflect.ValueOf(shadow.ToInt(in)))

	case reflect.Int8:
		r.setZero(out, path)
		out.Set(reflect.ValueOf(shadow.ToInt8(in)))

	case reflect.Int16:
		r.setZero(out, path)
		out.Set(reflect.ValueOf(shadow.ToInt16(in)))

	case reflect.Int32:
		r.setZero(out, path)
		out.Set(reflect.ValueOf(shadow.ToInt32(in)))

	case reflect.Int64:
		r.setZero(out, path)
		out.SetInt(shadow.ToInt64(in))

	case reflect.Float32:
		r.setZero(out, path)
		out.Set(reflect.ValueOf(shadow.ToFloat32(in)))

	case reflect.Float64:
		r.setZero(out, path)
		out.SetFloat(shadow.ToFloat64(in))

	default:
	}
}

func (r *Request) validateAny(out reflect.Value, path string) {
	out = reflect.Indirect(out)
	if !out.IsValid() {
		return
	}

	switch out.Kind() {
	case reflect.Struct:
		r.validateStruct(out, path)
	case reflect.Slice:
		for i := 0; i < out.Len(); i++ {
			r.validateAny(out.Index(i), r.getPath(path, fmt.Sprintf("[%d]", i)))
		}
	case reflect.Map:
		for _, n := range out.MapKeys() {
			r.validateAny(out.MapIndex(n), r.getPath(path, fmt.Sprintf("[%q]", n.String())))
		}
	}
}

func (r *Request) validateStruct(value reflect.Value, path string) {
	for i := 0; i < value.Type().NumField(); i++ {
		f := value.Type().Field(i)
		val := value.FieldByName(f.Name)

		name := f.Tag.Get(ParamNameTag)
		if name == "" {
			name = strings.ToLower(f.Name)
		}

		childPath := r.getPath(path, name)
		invalid := false

		if val.Kind() != reflect.Ptr {
			switch val.Kind() {
			case reflect.Map, reflect.Slice:
				if val.IsNil() {
					invalid = true
				}
			default:
				if !val.IsValid() {
					invalid = true
				}

				if !invalid && r.isZero(val) {
					_, ok := r.zeros[childPath]
					invalid = !ok
				}
			}
		}

		if invalid {
			r.errors = append(r.errors, "missing required parameter: "+childPath)
		} else {
			r.validateAny(val, childPath)
		}
	}
}
