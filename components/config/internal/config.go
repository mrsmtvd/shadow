package internal

import (
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(config.ConfigDebug, config.ValueTypeBool).
			WithUsage("Debug mode").
			WithGroup("Develop mode").
			WithEditable(true),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{config.WatcherForAll}, c.watchChanges),
	}
}

func (c *Component) watchChanges(key string, newValue interface{}, oldValue interface{}) {
	c.log().Infof("Change value for %s with '%v' to '%v'", key, oldValue, newValue)
}
