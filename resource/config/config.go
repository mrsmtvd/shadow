package config

const (
	ConfigDebug = "debug"
)

func (r *Resource) GetConfigVariables() []Variable {
	return []Variable{
		{
			Key:     ConfigDebug,
			Default: false,
			Usage:   "Debug mode",
			Type:    ValueTypeBool,
		},
	}
}
