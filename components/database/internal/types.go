package internal

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/go-gorp/gorp"
	"github.com/kihamo/gotypes"
)

type StructType map[string]interface{}

type TypeConverter struct{}

func (o StructType) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	for i := range o {
		data[gotypes.ToString(i)] = o[i]
	}

	return json.Marshal(data)
}

func (o *StructType) UnmarshalJSON(data []byte) error {
	target := make(map[string]interface{})

	if err := json.Unmarshal(data, &target); err != nil {
		return err
	}

	if target == nil {
		return &json.UnmarshalTypeError{Value: "null", Type: reflect.TypeOf(o)}
	}

	*o = StructType{}
	for i := range target {
		(*o)[i] = target[i]
	}

	return nil
}

func (t TypeConverter) ToDb(val interface{}) (interface{}, error) {
	switch t := val.(type) {
	case StructType:
		b, err := json.Marshal(t)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	return val, nil
}

func (c TypeConverter) FromDb(target interface{}) (gorp.CustomScanner, bool) {
	switch target.(type) {
	case *StructType:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New("FromDb: Unable to convert JsonField to *string")
			}
			return json.Unmarshal([]byte(*s), target)
		}
		return gorp.CustomScanner{new(string), target, binder}, true
	}

	return gorp.CustomScanner{}, false
}
