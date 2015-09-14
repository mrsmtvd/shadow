package api

import (
	"reflect"

	"github.com/kihamo/gotypes"
)

type Request struct {
	converter *gotypes.Converter
	args      []interface{}
	kwargs    map[string]interface{}
}

func NewRequest(output interface{}, args []interface{}, kwargs map[string]interface{}) *Request {
	request := &Request{
		args:   args,
		kwargs: kwargs,
	}

	switch reflect.Indirect(reflect.ValueOf(output)).Kind() {
	case reflect.Struct, reflect.Map, reflect.Interface:
		request.converter = gotypes.NewConverter(kwargs, output)

	default:
		request.converter = gotypes.NewConverter(args, output)
	}

	return request
}

func (r *Request) Valid() (bool, []string) {
	return r.converter.Valid(), r.converter.GetInvalidFields()
}

func (r *Request) GetArgs() []interface{} {
	return r.args
}

func (r *Request) GetKwargs() map[string]interface{} {
	return r.kwargs
}
