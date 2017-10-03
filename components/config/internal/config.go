package internal

import (
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			config.ConfigDebug,
			config.ValueTypeBool,
			false,
			"Debug mode",
			true,
			nil,
			nil),
	}
}

func (c *Component) GetConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher(config.ComponentName, []string{config.WatcherForAll}, c.watchChanges),
	}
}

func (c *Component) watchChanges(key string, newValue interface{}, oldValue interface{}) {
	c.log().Infof("Change value for %s with '%v' to '%v'", key, oldValue, newValue)
}
