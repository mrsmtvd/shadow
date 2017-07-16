package grpc

import (
	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigHost              = ComponentName + ".host"
	ConfigPort              = ComponentName + ".port"
	ConfigReflectionEnabled = ComponentName + ".reflection_enabled"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:     ConfigHost,
			Default: "localhost",
			Usage:   "gRPC host",
			Type:    config.ValueTypeString,
		},
		{
			Key:     ConfigPort,
			Default: 50052,
			Usage:   "gRPC port number",
			Type:    config.ValueTypeInt,
		},
		{
			Key:     ConfigReflectionEnabled,
			Default: true,
			Usage:   "gRPC enabled register reflection",
			Type:    config.ValueTypeBool,
		},
	}
}
