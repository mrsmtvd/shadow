package internal

import (
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/grpc"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			grpc.ConfigHost,
			config.ValueTypeString,
			"localhost",
			"gRPC host",
			false,
			"",
			nil,
			nil),
		config.NewVariable(
			grpc.ConfigPort,
			config.ValueTypeInt,
			50052,
			"gRPC port number",
			false,
			"",
			nil,
			nil),
		config.NewVariable(
			grpc.ConfigReflectionEnabled,
			config.ValueTypeInt,
			true,
			"gRPC enabled register reflection",
			false,
			"",
			nil,
			nil),
		config.NewVariable(
			grpc.ConfigManagerMaxLevel,
			config.ValueTypeInt,
			2,
			"Max level of parsing types",
			true,
			"",
			nil,
			nil),
	}
}
