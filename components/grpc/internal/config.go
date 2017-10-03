package internal

import (
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/grpc"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariableItem(
			grpc.ConfigHost,
			config.ValueTypeString,
			"localhost",
			"gRPC host",
			false,
			nil,
			nil),
		config.NewVariableItem(
			grpc.ConfigPort,
			config.ValueTypeInt,
			50052,
			"gRPC port number",
			false,
			nil,
			nil),
		config.NewVariableItem(
			grpc.ConfigReflectionEnabled,
			config.ValueTypeInt,
			true,
			"gRPC enabled register reflection",
			false,
			nil,
			nil),
	}
}
