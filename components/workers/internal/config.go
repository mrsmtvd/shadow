package internal

import (
	"time"

	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/logging"
	"github.com/mrsmtvd/shadow/components/workers"
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
		config.NewVariable(workers.ConfigListenersLoggingEnabled, config.ValueTypeBool).
			WithUsage("Logging listener enabled").
			WithEditable(true).
			WithDefault(true).
			WithGroup("listeners"),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{workers.ConfigWorkersCount}, c.watchCount),
		config.NewWatcher([]string{workers.ConfigTickerExecuteTasksDuration}, c.watchTickerExecuteTasksDuration),
		config.NewWatcher([]string{workers.ConfigListenersLoggingEnabled}, c.watchListenersLoggingEnabled),
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

func (c *Component) watchListenersLoggingEnabled(_ string, newValue interface{}, _ interface{}) {
	if newValue.(bool) {
		if l := c.newLoggingListener(); l != nil {
			c.addLockedListener(l)
		}

		return
	}

	c.removeLockedListener(c.Name() + "." + logging.ComponentName)
}
