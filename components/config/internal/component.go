package internal

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kihamo/gotypes"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logger"
)

func EnvKey(name string) string {
	name = strings.ToUpper(name)
	name = strings.Replace(name, " ", "_", -1)
	name = strings.Replace(name, "-", "_", -1)
	name = strings.Replace(name, ".", "_", -1)

	return name
}

type Component struct {
	mutex       sync.RWMutex
	application shadow.Application
	logger      logger.Logger
	envPrefix   string
	variables   map[string]config.Variable
	watchers    map[string][]config.Watcher
	routes      []dashboard.Route
}

func (c *Component) GetName() string {
	return config.ComponentName
}

func (c *Component) GetVersion() string {
	return config.ComponentVersion
}

func (c *Component) Init(a shadow.Application) (err error) {
	c.application = a
	c.envPrefix = EnvKey(a.GetName()) + "_"
	c.variables = make(map[string]config.Variable)
	c.watchers = make(map[string][]config.Watcher)

	return err
}

func (c *Component) Run() error {
	components, err := c.application.GetComponents()
	if err != nil {
		return err
	}

	for _, component := range components {
		if variables, ok := component.(config.HasVariables); ok {
			for _, variable := range variables.GetConfigVariables() {
				if variable.Key() == config.WatcherForAll {
					return fmt.Errorf("Use key %s not allowed", config.WatcherForAll)
				}

				c.mutex.Lock()
				c.variables[variable.Key()] = variable
				c.mutex.Unlock()
			}
		}
	}

	if err := c.LoadFromEnv(); err != nil {
		return err
	}

	if err := c.LoadFromCLIArguments(); err != nil {
		return err
	}

	for _, component := range components {
		if watchers, ok := component.(config.HasWatchers); ok {
			for _, watcher := range watchers.GetConfigWatchers() {
				c.Watch(watcher)
			}
		}
	}

	c.logConfig()

	return nil
}

func (c *Component) LoadFromCLIArguments() error {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	for _, v := range c.GetAllVariables() {
		switch v.Type() {
		case config.ValueTypeBool:
			flagSet.Bool(v.Key(), c.GetBool(v.Key()), v.Usage())
		case config.ValueTypeInt:
			flagSet.Int(v.Key(), c.GetInt(v.Key()), v.Usage())
		case config.ValueTypeInt64:
			flagSet.Int64(v.Key(), c.GetInt64(v.Key()), v.Usage())
		case config.ValueTypeUint:
			flagSet.Uint(v.Key(), c.GetUint(v.Key()), v.Usage())
		case config.ValueTypeUint64:
			flagSet.Uint64(v.Key(), c.GetUint64(v.Key()), v.Usage())
		case config.ValueTypeFloat64:
			flagSet.Float64(v.Key(), c.GetFloat64(v.Key()), v.Usage())
		case config.ValueTypeString:
			flagSet.String(v.Key(), c.GetString(v.Key()), v.Usage())
		case config.ValueTypeDuration:
			flagSet.Duration(v.Key(), c.GetDuration(v.Key()), v.Usage())
		}
	}

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return err
	}

	flagSet.Visit(func(f *flag.Flag) {
		if c.Has(f.Name) {
			c.Set(f.Name, f.Value.String())
		}
	})

	return nil
}

func (c *Component) LoadFromEnv() error {
	for _, v := range c.GetAllVariables() {
		envKey := c.envPrefix + EnvKey(v.Key())
		if value, ok := os.LookupEnv(envKey); ok {
			if err := c.Set(v.Key(), value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Component) EnvPrefix() string {
	return c.envPrefix
}

func (c *Component) GetWatchers(key string) []config.Watcher {
	watchers := []config.Watcher{}

	if watchersForAll, ok := c.watchers[config.WatcherForAll]; ok {
		for _, w := range watchersForAll {
			watchers = append(watchers, w)
		}
	}

	if key != config.WatcherForAll {
		if watchersByKey, ok := c.watchers[key]; ok {
			for _, w := range watchersByKey {
				watchers = append(watchers, w)
			}
		}
	}

	return watchers
}

func (c *Component) Watch(watcher config.Watcher) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, key := range watcher.Keys() {
		if watchers, ok := c.watchers[key]; ok {
			c.watchers[key] = append(watchers, watcher)
		} else {
			c.watchers[key] = []config.Watcher{watcher}
		}
	}
}

func (c *Component) log() logger.Logger {
	if c.logger == nil {
		c.logger = logger.NewOrNop(c.GetName(), c.application)
	}

	return c.logger
}

func (c *Component) logConfig() {
	fields := make(map[string]interface{}, 0)

	if c.envPrefix != "" {
		fields["config.prefix"] = c.envPrefix
	}

	for _, v := range c.GetAllVariables() {
		fields[v.Key()] = v.Value()
	}

	c.log().Info("Init config", fields)
}

func (c *Component) Has(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	_, ok := c.variables[key]
	return ok
}

func (c *Component) Get(key string) interface{} {
	c.mutex.RLock()
	v, ok := c.variables[key]
	c.mutex.RUnlock()

	if ok && v.Value() != nil {
		switch v.Type() {
		case config.ValueTypeBool:
			return c.GetBool(key)
		case config.ValueTypeInt:
			return c.GetInt(key)
		case config.ValueTypeInt64:
			return c.GetInt64(key)
		case config.ValueTypeUint:
			return c.GetUint(key)
		case config.ValueTypeUint64:
			return c.GetUint64(key)
		case config.ValueTypeFloat64:
			return c.GetFloat64(key)
		case config.ValueTypeString:
			return c.GetString(key)
		case config.ValueTypeDuration:
			return c.GetDuration(key)
		}
	}

	return nil
}

func (c *Component) Set(key string, value interface{}) error {
	if !c.Has(key) {
		return errors.New("Variable not found")
	}

	oldValue := c.Get(key)
	c.mutex.RLock()

	switch c.variables[key].Type() {
	case config.ValueTypeBool:
		value = gotypes.ToBool(value)
	case config.ValueTypeInt:
		value = gotypes.ToInt(value)
	case config.ValueTypeInt64:
		value = gotypes.ToInt64(value)
	case config.ValueTypeUint:
		value = gotypes.ToUint(value)
	case config.ValueTypeUint64:
		value = gotypes.ToUint64(value)
	case config.ValueTypeFloat64:
		value = gotypes.ToFloat64(value)
	case config.ValueTypeString:
		value = gotypes.ToString(value)
	case config.ValueTypeDuration:
		value = gotypes.ToDuration(value)
	default:
		c.mutex.RUnlock()
		return fmt.Errorf("Unknown type %s for config %s", c.variables[key].Type(), c.variables[key].Key())
	}

	err := c.variables[key].Change(value)
	c.mutex.RUnlock()

	if err != nil {
		return err
	}

	watchers := c.GetWatchers(key)

	if len(watchers) > 0 {
		go func() {
			for _, item := range watchers {
				item.Callback(key, value, oldValue)
			}
		}()
	}

	return nil
}

func (c *Component) IsEditable(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if variable, ok := c.variables[key]; ok {
		return variable.Editable()
	}

	return false
}

func (c *Component) GetAllVariables() map[string]config.Variable {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	variables := make(map[string]config.Variable, len(c.variables))

	for k, v := range c.variables {
		variables[k] = v
	}

	return variables
}

func (c *Component) GetBool(key string) bool {
	return c.GetBoolDefault(key, false)
}

func (c *Component) GetBoolDefault(key string, value interface{}) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToBool(val.Value())
	}

	return gotypes.ToBool(value)
}

func (c *Component) GetInt(key string) int {
	return c.GetIntDefault(key, -1)
}

func (c *Component) GetIntDefault(key string, value interface{}) int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToInt(val.Value())
	}

	return gotypes.ToInt(value)
}

func (c *Component) GetInt64(key string) int64 {
	return c.GetInt64Default(key, -1)
}

func (c *Component) GetInt64Default(key string, value interface{}) int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToInt64(val.Value())
	}

	return gotypes.ToInt64(value)
}

func (c *Component) GetUint(key string) uint {
	return c.GetUintDefault(key, 0)
}

func (c *Component) GetUintDefault(key string, value interface{}) uint {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToUint(val.Value())
	}

	return gotypes.ToUint(value)
}

func (c *Component) GetUint64(key string) uint64 {
	return c.GetUint64Default(key, 0)
}

func (c *Component) GetUint64Default(key string, value interface{}) uint64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToUint64(val.Value())
	}

	return gotypes.ToUint64(value)
}

func (c *Component) GetFloat64(key string) float64 {
	return c.GetFloat64Default(key, -1)
}

func (c *Component) GetFloat64Default(key string, value interface{}) float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToFloat64(val.Value())
	}

	return gotypes.ToFloat64(value)
}

func (c *Component) GetString(key string) string {
	return c.GetStringDefault(key, "")
}

func (c *Component) GetStringDefault(key string, value interface{}) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToString(val.Value())
	}

	return gotypes.ToString(value)
}

func (c *Component) GetDuration(key string) time.Duration {
	return c.GetDurationDefault(key, 0)
}

func (c *Component) GetDurationDefault(key string, value interface{}) time.Duration {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToDuration(val.Value())
	}

	return gotypes.ToDuration(value)
}
