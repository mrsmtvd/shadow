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
	Group() string
	View() []string
	ViewOptions() map[string]interface{}
	Change(value interface{}) error
}

type HasVariables interface {
	ConfigVariables() []Variable
}

type variableValue struct {
	original interface{}
}

type VariableSimple struct {
	key         string
	typ         string
	def         interface{}
	value       atomic.Value
	usage       string
	editable    bool
	group       string
	view        []string
	viewOptions map[string]interface{}
}

func NewVariable(key string, typ string, def interface{}, usage string, editable bool, group string, view []string, viewOptions map[string]interface{}) *VariableSimple {
	v := &VariableSimple{
		key:         key,
		typ:         typ,
		def:         def,
		usage:       usage,
		editable:    editable,
		group:       group,
		view:        view,
		viewOptions: viewOptions,
	}
	v.Change(def)

	return v
}

func (v *VariableSimple) Key() string {
	return v.key
}

func (v *VariableSimple) Type() string {
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

func (v *VariableSimple) Default() interface{} {
	return v.def
}

func (v *VariableSimple) Value() interface{} {
	var value interface{}

	if l := v.value.Load(); l != nil {
		value = l.(*variableValue).original
	}

	return value
}

func (v *VariableSimple) Usage() string {
	return v.usage
}

func (v *VariableSimple) Editable() bool {
	return v.editable
}

func (v *VariableSimple) Group() string {
	return v.group
}

func (v *VariableSimple) View() []string {
	return v.view
}

func (v *VariableSimple) ViewOptions() map[string]interface{} {
	return v.viewOptions
}

func (v *VariableSimple) Change(value interface{}) error {
	v.value.Store(&variableValue{value})
	return nil
}

func (v *VariableSimple) String() string {
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

func (v *VariableSimple) GoString() string {
	return v.String()
}
