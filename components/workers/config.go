package workers

import (
	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigWorkersCount = "workers.count"
	ConfigDoneSize     = "workers.done.size"
)

func (c *Component) GetConfigVariables() []config.Variable {
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

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigWorkersCount: {c.watchWorkersCount},
	}
}

func (c *Component) watchWorkersCount(newValue interface{}, _ interface{}) {
	for i := c.dispatcher.GetWorkers().Len(); i < newValue.(int); i++ {
		c.AddWorker()
	}
}
