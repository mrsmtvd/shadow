package internal

import (
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/workers"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			workers.ConfigWorkersCount,
			config.ValueTypeInt,
			2,
			"Default workers count",
			true,
			nil,
			nil),
	}
}

func (c *Component) GetConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher(workers.ComponentName, []string{workers.ConfigWorkersCount}, c.watchCount),
	}
}

func (c *Component) watchCount(_ string, newValue interface{}, _ interface{}) {
	for i := len(c.GetWorkers()); i < newValue.(int); i++ {
		c.AddSimpleWorker()
	}
}
