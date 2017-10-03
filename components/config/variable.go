package config

import (
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

type VariableItem struct {
	key         string
	typ         string
	def         interface{}
	value       interface{}
	usage       string
	editable    bool
	view        []string
	viewOptions map[string]interface{}
}

func NewVariable(key string, typ string, def interface{}, usage string, editable bool, view []string, viewOptions map[string]interface{}) Variable {
	return &VariableItem{
		key:         key,
		typ:         typ,
		def:         def,
		value:       def,
		usage:       usage,
		editable:    editable,
		view:        view,
		viewOptions: viewOptions,
	}
}

func (v *VariableItem) Key() string {
	return v.key
}

func (v *VariableItem) Type() string {
	// autodetect type of value
	if v.typ == "" && (v.def != nil || v.value != nil) {
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
	return v.value
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
	v.value = value
	return nil
}
