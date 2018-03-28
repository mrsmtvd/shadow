package internal

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
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

type variableSort struct {
	variable config.Variable
	order    int
}

type Component struct {
	mutex         sync.RWMutex
	application   shadow.Application
	logger        logger.Logger
	envPrefix     string
	variables     map[string]config.Variable
	variablesSort []*variableSort
	watchers      map[string][]*WatcherItem
	routes        []dashboard.Route
}

func (c *Component) Name() string {
	return config.ComponentName
}

func (c *Component) Version() string {
	return config.ComponentVersion
}

func (c *Component) Init(a shadow.Application) (err error) {
	components, err := a.GetComponents()
	if err != nil {
		return err
	}

	c.application = a
	c.envPrefix = EnvKey(a.Name()) + "_"
	c.variables = make(map[string]config.Variable)
	c.variablesSort = make([]*variableSort, 0)
	c.watchers = make(map[string][]*WatcherItem)

	for _, component := range components {
		if cmpVariables, ok := component.(config.HasVariables); ok {
			for _, variable := range cmpVariables.ConfigVariables() {
				if variable.Key() == config.WatcherForAll {
					return fmt.Errorf("Use key %s not allowed", config.WatcherForAll)
				}

				variablesForSort := &variableSort{
					variable: variable,
				}

				c.mutex.Lock()
				c.variables[variable.Key()] = variable

				variablesForSort.order = len(c.variables)
				c.variablesSort = append(c.variablesSort, variablesForSort)
				c.mutex.Unlock()
			}
		}
	}

	c.mutex.Lock()
	sort.SliceStable(c.variablesSort, func(i, j int) bool {
		return c.variablesSort[i].order < c.variablesSort[j].order
	})
	c.mutex.Unlock()

	if err := c.LoadFromEnv(); err != nil {
		return err
	}

	if err := c.LoadFromCLIArguments(); err != nil {
		return err
	}

	for _, component := range components {
		if watchers, ok := component.(config.HasWatchers); ok {
			for _, watcher := range watchers.ConfigWatchers() {
				c.Watch(watcher, component.Name())
			}
		}
	}

	return err
}

func (c *Component) Run() error {
	fields := make(map[string]interface{}, 0)

	if c.envPrefix != "" {
		fields["config.prefix"] = c.envPrefix
	}

	for _, v := range c.Variables() {
		fields[v.Key()] = v.Value()
	}

	c.log().Info("Init config", fields)

	return nil
}

func (c *Component) LoadFromCLIArguments() error {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	for _, v := range c.Variables() {
		switch v.Type() {
		case config.ValueTypeBool:
			flagSet.Bool(v.Key(), c.Bool(v.Key()), v.Usage())
		case config.ValueTypeInt:
			flagSet.Int(v.Key(), c.Int(v.Key()), v.Usage())
		case config.ValueTypeInt64:
			flagSet.Int64(v.Key(), c.Int64(v.Key()), v.Usage())
		case config.ValueTypeUint:
			flagSet.Uint(v.Key(), c.Uint(v.Key()), v.Usage())
		case config.ValueTypeUint64:
			flagSet.Uint64(v.Key(), c.Uint64(v.Key()), v.Usage())
		case config.ValueTypeFloat64:
			flagSet.Float64(v.Key(), c.Float64(v.Key()), v.Usage())
		case config.ValueTypeString:
			flagSet.String(v.Key(), c.String(v.Key()), v.Usage())
		case config.ValueTypeDuration:
			flagSet.Duration(v.Key(), c.Duration(v.Key()), v.Usage())
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
	for _, v := range c.Variables() {
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

func (c *Component) Watchers(key string) []config.Watcher {
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

func (c *Component) Watch(watcher config.Watcher, source string) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item := NewWatcherItem(watcher, source)

	for _, key := range watcher.Keys() {
		if watchers, ok := c.watchers[key]; ok {
			c.watchers[key] = append(watchers, item)
		} else {
			c.watchers[key] = []*WatcherItem{item}
		}
	}
}

func (c *Component) log() logger.Logger {
	if c.logger == nil {
		c.logger = logger.NewOrNop(c.Name(), c.application)
	}

	return c.logger
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
			return c.Bool(key)
		case config.ValueTypeInt:
			return c.Int(key)
		case config.ValueTypeInt64:
			return c.Int64(key)
		case config.ValueTypeUint:
			return c.Uint(key)
		case config.ValueTypeUint64:
			return c.Uint64(key)
		case config.ValueTypeFloat64:
			return c.Float64(key)
		case config.ValueTypeString:
			return c.String(key)
		case config.ValueTypeDuration:
			return c.Duration(key)
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

	watchers := c.Watchers(key)

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

func (c *Component) Variables() []config.Variable {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	variables := make([]config.Variable, 0, len(c.variablesSort))

	for _, v := range c.variablesSort {
		variables = append(variables, v.variable)
	}

	return variables
}

func (c *Component) Bool(key string) bool {
	return c.BoolDefault(key, false)
}

func (c *Component) BoolDefault(key string, value interface{}) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToBool(val.Value())
	}

	return gotypes.ToBool(value)
}

func (c *Component) Int(key string) int {
	return c.IntDefault(key, -1)
}

func (c *Component) IntDefault(key string, value interface{}) int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToInt(val.Value())
	}

	return gotypes.ToInt(value)
}

func (c *Component) Int64(key string) int64 {
	return c.Int64Default(key, -1)
}

func (c *Component) Int64Default(key string, value interface{}) int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToInt64(val.Value())
	}

	return gotypes.ToInt64(value)
}

func (c *Component) Uint(key string) uint {
	return c.UintDefault(key, 0)
}

func (c *Component) UintDefault(key string, value interface{}) uint {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToUint(val.Value())
	}

	return gotypes.ToUint(value)
}

func (c *Component) Uint64(key string) uint64 {
	return c.Uint64Default(key, 0)
}

func (c *Component) Uint64Default(key string, value interface{}) uint64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToUint64(val.Value())
	}

	return gotypes.ToUint64(value)
}

func (c *Component) Float64(key string) float64 {
	return c.Float64Default(key, -1)
}

func (c *Component) Float64Default(key string, value interface{}) float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToFloat64(val.Value())
	}

	return gotypes.ToFloat64(value)
}

func (c *Component) String(key string) string {
	return c.StringDefault(key, "")
}

func (c *Component) StringDefault(key string, value interface{}) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToString(val.Value())
	}

	return gotypes.ToString(value)
}

func (c *Component) Duration(key string) time.Duration {
	return c.DurationDefault(key, 0)
}

func (c *Component) DurationDefault(key string, value interface{}) time.Duration {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToDuration(val.Value())
	}

	return gotypes.ToDuration(value)
}
