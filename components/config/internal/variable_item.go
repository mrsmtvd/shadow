package internal

import (
	"github.com/mrsmtvd/shadow/components/config"
)

type VariableItem struct {
	config.Variable

	source   string
	variable config.Variable
}

func NewVariableItem(variable config.Variable, source string) *VariableItem {
	if source == "" {
		source = "unknown"
	}

	return &VariableItem{
		variable: variable,
		source:   source,
	}
}

func (v *VariableItem) Source() string {
	return v.source
}

func (v *VariableItem) Key() string {
	return v.variable.Key()
}

func (v *VariableItem) Default() interface{} {
	return v.variable.Default()
}

func (v *VariableItem) Value() interface{} {
	return v.variable.Value()
}

func (v *VariableItem) Type() string {
	return v.variable.Type()
}

func (v *VariableItem) Usage() string {
	return v.variable.Usage()
}

func (v *VariableItem) Editable() bool {
	return v.variable.Editable()
}

func (v *VariableItem) Group() string {
	return v.variable.Group()
}

func (v *VariableItem) View() []string {
	return v.variable.View()
}

func (v *VariableItem) ViewOptions() map[string]interface{} {
	return v.variable.ViewOptions()
}

func (v *VariableItem) Change(value interface{}) error {
	return v.variable.Change(value)
}

func (v *VariableItem) String() string {
	return v.variable.String()
}

func (v *VariableItem) GoString() string {
	return v.variable.GoString()
}
