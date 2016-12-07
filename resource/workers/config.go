package workers

import (
	"github.com/kihamo/shadow/resource/config"
)

func (r *Resource) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:   "workers.count",
			Value: 2,
			Usage: "Default workers count",
		},
		{
			Key:   "workers.done.size",
			Value: 1000,
			Usage: "Size buffer of done task channel",
		},
	}
}
