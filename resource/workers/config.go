package workers

import (
	"github.com/kihamo/shadow/resource/config"
)

const (
	ConfigWorkersCount = "workers.count"
	ConfigDoneSize     = "workers.done.size"
)

func (r *Resource) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:      ConfigWorkersCount,
			Default:  2,
			Usage:    "Default workers count",
			Type:     config.ValueTypeInt,
			Editable: true,
		},
		{
			Key:     ConfigDoneSize,
			Default: 1000,
			Usage:   "Size buffer of done task channel",
			Type:    config.ValueTypeInt,
		},
	}
}

func (r *Resource) GetConfigWatcher() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigWorkersCount: {
			r.watchWorkersCount,
		},
	}
}

func (r *Resource) watchWorkersCount(newValue interface{}, lastValue interface{}) {
	for i := r.dispatcher.GetWorkers().Len(); i < r.config.GetInt(ConfigWorkersCount); i++ {
		r.dispatcher.AddWorker()
	}
}
