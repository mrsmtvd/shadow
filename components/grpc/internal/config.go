package internal

import (
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/grpc"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			grpc.ConfigHost,
			config.ValueTypeString,
			"localhost",
			"Host",
			false,
			"Lister",
			nil,
			nil),
		config.NewVariable(
			grpc.ConfigPort,
			config.ValueTypeInt,
			50052,
			"Port number",
			false,
			"Lister",
			nil,
			nil),
		config.NewVariable(
			grpc.ConfigReflectionEnabled,
			config.ValueTypeInt,
			true,
			"Enabled register reflection",
			false,
			"Others",
			nil,
			nil),
		config.NewVariable(
			grpc.ConfigManagerMaxLevel,
			config.ValueTypeInt,
			2,
			"Max level of parsing types",
			true,
			"Others",
			nil,
			nil),
	}
}
