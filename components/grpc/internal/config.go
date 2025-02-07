package internal

import (
	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/grpc"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(grpc.ConfigHost, config.ValueTypeString).
			WithUsage("Host").
			WithGroup("Lister").
			WithDefault("localhost"),
		config.NewVariable(grpc.ConfigPort, config.ValueTypeString).
			WithUsage("Port number").
			WithGroup("Lister").
			WithDefault(50052),
		config.NewVariable(grpc.ConfigReflectionEnabled, config.ValueTypeBool).
			WithUsage("Enabled register reflection"),
		config.NewVariable(grpc.ConfigManagerMaxLevel, config.ValueTypeInt).
			WithUsage("Max level of parsing types").
			WithEditable(true).
			WithDefault(2),
	}
}
