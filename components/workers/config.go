package workers

import (
	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigWorkersCount = "workers.count"
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
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigWorkersCount: {c.watchWorkersCount},
	}
}

func (c *Component) watchWorkersCount(_ string, newValue interface{}, _ interface{}) {
	for i := c.dispatcher.GetWorkers().Len(); i < newValue.(int); i++ {
		c.AddWorker()
	}
}
