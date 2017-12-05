package config

import (
	"fmt"
	"sync/atomic"
	"time"
)

const (
	ValueTypeBool     = "bool"
	ValueTypeInt      = "int"
	ValueTypeInt64    = "int64"
	ValueTypeUint     = "uint"
	ValueTypeUint64   = "uint64"
	ValueTypeFloat64  = "float64"
	ValueTypeString   = "string"
	ValueTypeDuration = "duration"
)

type Variable interface {
	fmt.Stringer
	fmt.GoStringer

	Key() string
	Default() interface{}
	Value() interface{}
	Type() string
	Usage() string
	Editable() bool
	View() []string
	ViewOptions() map[string]interface{}
	Change(value interface{}) error
}

type HasVariables interface {
	GetConfigVariables() []Variable
}

type variableValue struct {
	original interface{}
}

type VariableItem struct {
	key         string
	typ         string
	def         interface{}
	value       atomic.Value
	usage       string
	editable    bool
	view        []string
	viewOptions map[string]interface{}
}

func NewVariable(key string, typ string, def interface{}, usage string, editable bool, view []string, viewOptions map[string]interface{}) Variable {
	v := &VariableItem{
		key:         key,
		typ:         typ,
		def:         def,
		usage:       usage,
		editable:    editable,
		view:        view,
		viewOptions: viewOptions,
	}
	v.Change(def)

	return v
}

func (v *VariableItem) Key() string {
	return v.key
}

func (v *VariableItem) Type() string {
	// autodetect type of value
	if v.typ == "" && (v.Default() != nil || v.Value() != nil) {
		switch v.Value().(type) {
		case bool:
			v.typ = ValueTypeBool
		case int:
			v.typ = ValueTypeInt
		case int64:
			v.typ = ValueTypeInt64
		case uint:
			v.typ = ValueTypeUint
		case uint64:
			v.typ = ValueTypeUint64
		case float64:
			v.typ = ValueTypeFloat64
		case string:
			v.typ = ValueTypeString
		case time.Duration:
			v.typ = ValueTypeDuration
		}
	}

	return v.typ
}

func (v *VariableItem) Default() interface{} {
	return v.def
}

func (v *VariableItem) Value() interface{} {
	var value interface{}

	if l := v.value.Load(); l != nil {
		value = l.(*variableValue).original
	}

	return value
}

func (v *VariableItem) Usage() string {
	return v.usage
}

func (v *VariableItem) Editable() bool {
	return v.editable
}

func (v *VariableItem) View() []string {
	return v.view
}

func (v *VariableItem) ViewOptions() map[string]interface{} {
	return v.viewOptions
}

func (v *VariableItem) Change(value interface{}) error {
	v.value.Store(&variableValue{value})
	return nil
}

func (v *VariableItem) String() string {
	value := v.Value()
	viewOptions := v.ViewOptions()

	for _, view := range v.View() {
		if view != ViewEnum {
			continue
		}

		opts, ok := viewOptions[ViewOptionEnumOptions]
		if !ok {
			continue
		}

		sliceOpts, ok := opts.([][]interface{})
		if !ok {
			continue
		}

		for _, optValue := range sliceOpts {
			if len(optValue) > 1 && optValue[0] == value {
				return fmt.Sprintf("%s", optValue[1])
			}
		}
	}

	return fmt.Sprintf("%s", value)
}

func (v *VariableItem) GoString() string {
	return v.String()
}
