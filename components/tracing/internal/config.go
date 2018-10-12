package internal

import (
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/tracing"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(tracing.ConfigAgentHost, config.ValueTypeString).
			WithUsage("Host").
			WithGroup("Agent").
			WithDefault("localhost"),
		config.NewVariable(tracing.ConfigAgentPort, config.ValueTypeInt).
			WithUsage("Port number").
			WithGroup("Agent").
			WithDefault(6831),
	}
}
