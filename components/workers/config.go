package workers

import (
	"time"

	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigCount                         = ComponentName + ".count"
	ConfigTickerExecuteTasksDuration    = ComponentName + ".ticker_execute_tasks_duration"
	ConfigTickerNotifyListenersDuration = ComponentName + ".ticker_notify_listeners_duration"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:      ConfigCount,
			Default:  2,
			Usage:    "Default workers count",
			Type:     config.ValueTypeInt,
			Editable: true,
		},
		{
			Key:      ConfigTickerExecuteTasksDuration,
			Default:  "1s",
			Usage:    "Duration for ticker for execute tasks",
			Type:     config.ValueTypeDuration,
			Editable: true,
		},
		{
			Key:      ConfigTickerNotifyListenersDuration,
			Default:  "1s",
			Usage:    "Duration for notify listeners",
			Type:     config.ValueTypeDuration,
			Editable: true,
		},
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigCount:                         {c.watchCount},
		ConfigTickerExecuteTasksDuration:    {c.watchTickerExecuteTasksDuration},
		ConfigTickerNotifyListenersDuration: {c.watchTickerNotifyListenersDuration},
	}
}

func (c *Component) watchCount(_ string, newValue interface{}, _ interface{}) {
	for i := c.dispatcher.GetWorkers().Len(); i < newValue.(int); i++ {
		c.AddWorker()
	}
}

func (c *Component) watchTickerExecuteTasksDuration(_ string, newValue interface{}, _ interface{}) {
	c.dispatcher.SetTickerExecuteTasksDuration(newValue.(time.Duration))
}

func (c *Component) watchTickerNotifyListenersDuration(_ string, newValue interface{}, _ interface{}) {
	c.dispatcher.SetTickerNotifyListenersDuration(newValue.(time.Duration))
}
