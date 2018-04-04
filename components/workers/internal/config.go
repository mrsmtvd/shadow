package internal

import (
	"time"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/workers"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(workers.ConfigWorkersCount, config.ValueTypeInt).
			WithUsage("Default workers count").
			WithEditable(true).
			WithDefault(2),
		config.NewVariable(workers.ConfigTickerExecuteTasksDuration, config.ValueTypeDuration).
			WithUsage("Duration for ticker in dispatcher of workers").
			WithEditable(true).
			WithDefault("1s"),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{workers.ConfigWorkersCount}, c.watchCount),
	}
}

func (c *Component) watchCount(_ string, newValue interface{}, _ interface{}) {
	for i := len(c.GetWorkers()); i < newValue.(int); i++ {
		c.AddSimpleWorker()
	}
}

func (c *Component) watchTickerExecuteTasksDuration(_ string, newValue interface{}, _ interface{}) {
	c.dispatcher.SetTickerExecuteTasksDuration(newValue.(time.Duration))
}
