package internal

import (
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariableItem(
			config.ConfigDebug,
			config.ValueTypeBool,
			false,
			"Debug mode",
			true,
			nil,
			nil),
	}
}
