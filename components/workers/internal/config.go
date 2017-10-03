package internal

import (
	"time"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/workers"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariableItem(
			workers.ConfigCount,
			config.ValueTypeInt,
			2,
			"Default workers count",
			true,
			nil,
			nil),
		config.NewVariableItem(
			workers.ConfigTickerExecuteTasksDuration,
			config.ValueTypeDuration,
			"1s",
			"Duration for ticker for execute tasks",
			true,
			nil,
			nil),
		config.NewVariableItem(
			workers.ConfigTickerNotifyListenersDuration,
			config.ValueTypeDuration,
			"1s",
			"Duration for notify listeners",
			true,
			nil,
			nil),
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		workers.ConfigCount:                         {c.watchCount},
		workers.ConfigTickerExecuteTasksDuration:    {c.watchTickerExecuteTasksDuration},
		workers.ConfigTickerNotifyListenersDuration: {c.watchTickerNotifyListenersDuration},
	}
}

func (c *Component) watchCount(_ string, newValue interface{}, _ interface{}) {
	for i := len(c.dispatcher.GetWorkers()); i < newValue.(int); i++ {
		c.AddWorker()
	}
}

func (c *Component) watchTickerExecuteTasksDuration(_ string, newValue interface{}, _ interface{}) {
	c.dispatcher.SetTickerExecuteTasksDuration(newValue.(time.Duration))
}

func (c *Component) watchTickerNotifyListenersDuration(_ string, newValue interface{}, _ interface{}) {
	c.dispatcher.SetTickerNotifyListenersDuration(newValue.(time.Duration))
}
