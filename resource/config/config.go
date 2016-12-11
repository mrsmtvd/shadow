package config

func (r *Resource) GetConfigVariables() []Variable {
	return []Variable{
		{
			Key:     "debug",
			Default: false,
			Usage:   "Debug mode",
			Type:    ValueTypeBool,
		},
	}
}
