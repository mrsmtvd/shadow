package api

import (
	"reflect"
	"strconv"
	"strings"
)

const (
	ParamNameTag = "param"
)

func RequestFillAndValidate(out interface{}, inArgs []interface{}, inKwargs map[string]interface{}) error {
	RequestFill(out, inArgs, inKwargs)

	if err := RequestValidate(out); err != nil {
		return err
	}

	return nil
}

func RequestFill(out interface{}, inArgs []interface{}, inKwargs map[string]interface{}) error {
	outValue := reflect.ValueOf(&out)

	if outValue.Kind() == reflect.Ptr && outValue.IsValid() {
		outValue = outValue.Elem()
	}

	switch outValue.Kind() {
	case reflect.Struct, reflect.Map, reflect.Interface:
		fill(outValue, inKwargs)

	default:
		fill(outValue, inArgs)
	}

	return nil
}

func RequestValidate(out interface{}) error {
	return nil
}

func fill(out reflect.Value, in interface{}) {
	switch out.Kind() {
	case reflect.Ptr:
		if out.IsValid() {
			fill(out.Elem(), in)
		}

	case reflect.Interface:
		fill(out.Elem(), in)

	case reflect.Struct:
		if in, ok := in.(map[string]interface{}); ok {
			for i := 0; i < out.NumField(); i++ {
				f := out.Type().Field(i)

				name := f.Tag.Get(ParamNameTag)
				if name == "" {
					name = strings.ToLower(f.Name)
				}

				if value, ok := in[name]; ok {
					fill(out.Field(i), value)
				}
			}
		}

	case reflect.Map:
		if in, ok := in.(map[string]interface{}); ok {
			keyType := out.Type().Key()
			valueType := out.Type().Elem()

			out.Set(reflect.MakeMap(out.Type()))
			for mapKey, mapValue := range in {
				key := reflect.New(keyType).Elem()
				fill(key, mapKey)

				value := reflect.New(valueType).Elem()
				fill(value, mapValue)

				out.SetMapIndex(key, value)
			}
		}

	case reflect.Slice:
		if in, ok := in.([]interface{}); ok {
			out.Set(reflect.MakeSlice(out.Type(), len(in), cap(in)))
			for i := range in {
				fill(out.Index(i), in[i])
			}
		}

	case reflect.String:
		var castIn string

		switch v := in.(type) {
		case int:
			castIn = strconv.Itoa(v)

		case float64:
			castIn = strconv.FormatFloat(v, 'f', 6, 64)

		default:
			if cast, ok := in.(string); ok {
				castIn = cast
			}
		}

		out.Set(reflect.ValueOf(castIn))

	case reflect.Bool:
		var castIn bool

		switch v := in.(type) {
		case bool:
			castIn = v
		case string:
			castIn = v != "" && v != "0" && v != "false"
		case int, float64, int64, uint64:
			castIn = v != 0
		}

		out.Set(reflect.ValueOf(castIn))

	case reflect.Float64:
		var castIn float64

		switch v := in.(type) {
		case string:
			castIn, _ = strconv.ParseFloat(v, 64)

		case float64:
			castIn = v

		case int:
			castIn = float64(v)

		case int64:
			castIn = float64(v)

		case uint:
			castIn = float64(v)

		case uint64:
			castIn = float64(v)
		}

		out.Set(reflect.ValueOf(castIn))

	case reflect.Int:
		var castIn int

		switch v := in.(type) {
		case string:
			castIn, _ = strconv.Atoi(v)

		case int:
			castIn = v

		case int64:
			castIn = int(v)

		case float64:
			castIn = int(v)

		case uint:
			castIn = int(v)

		case uint64:
			castIn = int(v)
		}

		out.Set(reflect.ValueOf(castIn))

	case reflect.Int64:
		var castIn int64

		switch v := in.(type) {
		case string:
			t, _ := strconv.Atoi(v)
			castIn = int64(t)

		case int64:
			castIn = v

		case int:
			castIn = int64(v)

		case float64:
			castIn = int64(v)

		case uint:
			castIn = int64(v)

		case uint64:
			castIn = int64(v)
		}

		out.Set(reflect.ValueOf(castIn))

	default:
		out.Set(reflect.ValueOf(in))
	}
}

func validate(out reflect.Value) error {
	return nil
}
