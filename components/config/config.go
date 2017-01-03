package config

const (
	ConfigDebug = "debug"
)

func (c *Component) GetConfigVariables() []Variable {
	return []Variable{
		{
			Key:      ConfigDebug,
			Default:  false,
			Usage:    "Debug mode",
			Type:     ValueTypeBool,
			Editable: true,
		},
	}
}
