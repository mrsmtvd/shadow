package config

const (
	ConfigDebug = ComponentName + ".debug"
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
