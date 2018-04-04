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
	def         atomic.Value
	value       atomic.Value
	usage       string
	editable    bool
	group       string
	view        []string
	viewOptions map[string]interface{}
}

func NewVariable(key string, typ string) *VariableSimple {
	return &VariableSimple{
		key:         key,
		typ:         typ,
		group:       "Others",
		view:        make([]string, 0, 0),
		viewOptions: make(map[string]interface{}, 0),
	}
}

func (v *VariableSimple) Key() string {
	return v.key
}
func (v *VariableSimple) WithKey(key string) *VariableSimple {
	v.key = key
	return v
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

func (v *VariableSimple) WithType(typ string) *VariableSimple {
	switch v.Type() {
	case ValueTypeInt, ValueTypeInt64, ValueTypeUint, ValueTypeUint64, ValueTypeFloat64, ValueTypeBool, ValueTypeString, ValueTypeDuration:
		v.typ = typ
	default:
		panic("Unknown type " + typ)
	}

	return v
}

func (v *VariableSimple) Default() interface{} {
	if l := v.def.Load(); l != nil {
		return l.(*variableValue).original
	}

	return nil
}

func (v *VariableSimple) WithDefault(value interface{}) *VariableSimple {
	v.def.Store(&variableValue{value})

	if l := v.value.Load(); l == nil {
		v.Change(value)
	}

	return v
}

func (v *VariableSimple) Value() interface{} {
	if l := v.value.Load(); l != nil {
		return l.(*variableValue).original
	}

	return nil
}

func (v *VariableSimple) WithValue(value interface{}) *VariableSimple {
	v.Change(value)
	return v
}

func (v *VariableSimple) Usage() string {
	return v.usage
}

func (v *VariableSimple) WithUsage(usage string) *VariableSimple {
	v.usage = usage
	return v
}

func (v *VariableSimple) Editable() bool {
	return v.editable
}

func (v *VariableSimple) WithEditable(editable bool) *VariableSimple {
	v.editable = editable
	return v
}

func (v *VariableSimple) Group() string {
	return v.group
}

func (v *VariableSimple) WithGroup(group string) *VariableSimple {
	v.group = group
	return v
}

func (v *VariableSimple) View() []string {
	return v.view
}

func (v *VariableSimple) WithView(view []string) *VariableSimple {
	v.view = view
	return v
}

func (v *VariableSimple) ViewOptions() map[string]interface{} {
	return v.viewOptions
}

func (v *VariableSimple) WithViewOptions(viewOptions map[string]interface{}) *VariableSimple {
	v.viewOptions = viewOptions
	return v
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
