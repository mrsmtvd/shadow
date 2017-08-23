package config

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/kihamo/gotypes"
	"github.com/kihamo/shadow"
	"github.com/rakyll/globalconf"
)

const (
	ComponentName = "config"

	FlagConfig = "config"

	WatcherForAll = "*"

	ValueTypeBool     = "bool"
	ValueTypeInt      = "int"
	ValueTypeInt64    = "int64"
	ValueTypeUint     = "uint"
	ValueTypeUint64   = "uint64"
	ValueTypeFloat64  = "float64"
	ValueTypeString   = "string"
	ValueTypeDuration = "duration"
)

type Watcher func(string, interface{}, interface{})

type Component struct {
	mutex       sync.RWMutex
	application shadow.Application
	config      *globalconf.GlobalConf
	variables   map[string]Variable
	watchers    map[string][]Watcher
}

type hasVariables interface {
	GetConfigVariables() []Variable
}

type hasWatchers interface {
	GetConfigWatchers() map[string][]Watcher
}

func (c *Component) GetName() string {
	return ComponentName
}

func (c *Component) GetVersion() string {
	return ComponentVersion
}

func (c *Component) Init(a shadow.Application) (err error) {
	c.application = a

	config := flag.String(FlagConfig, "", "Config file which which override default config parameters")
	flag.Parse()

	opts := globalconf.Options{
		EnvPrefix: strings.ToUpper(strings.Replace(a.GetName(), " ", "_", -1)) + "_",
		Filename:  *config,
	}

	if c.config, err = globalconf.NewWithOptions(&opts); err != nil {
		return err
	}

	c.mutex.Lock()
	c.variables = map[string]Variable{}
	c.watchers = map[string][]Watcher{}
	c.mutex.Unlock()

	return err
}

func (c *Component) Run() error {
	components, err := c.application.GetComponents()
	if err != nil {
		return err
	}

	for _, component := range components {
		if variables, ok := component.(hasVariables); ok {
			for _, variable := range variables.GetConfigVariables() {
				if variable.Key == WatcherForAll {
					return fmt.Errorf("Use key %s not allowed", WatcherForAll)
				}

				c.addFlag(variable)
			}
		}

		if watchers, ok := component.(hasWatchers); ok {
			for key, list := range watchers.GetConfigWatchers() {
				for _, watcher := range list {
					c.WatchVariable(key, watcher)
				}
			}
		}
	}

	c.config.ParseAll()

	c.mutex.Lock()
	for k, v := range c.variables {
		if reflect.ValueOf(v.Value).Kind() == reflect.Ptr {
			v.Value = reflect.Indirect(reflect.ValueOf(v.Value)).Interface()
			c.variables[k] = v
		}
	}
	c.mutex.Unlock()

	return nil
}

func (c *Component) addFlag(v Variable) {
	// autodetect type of value
	if v.Type == "" && (v.Default != nil || v.Value != nil) {
		var baseType interface{}

		if v.Default != nil {
			baseType = v.Default
		} else {
			baseType = v.Value
		}

		switch baseType.(type) {
		case bool:
			v.Type = ValueTypeBool
		case int:
			v.Type = ValueTypeInt
		case int64:
			v.Type = ValueTypeInt64
		case uint:
			v.Type = ValueTypeUint
		case uint64:
			v.Type = ValueTypeUint64
		case float64:
			v.Type = ValueTypeFloat64
		case string:
			v.Type = ValueTypeString
		case time.Duration:
			v.Type = ValueTypeDuration
		}
	}

	if v.Value == nil {
		v.Value = v.Default
	}

	switch v.Type {
	case ValueTypeBool:
		v.Value = flag.Bool(v.Key, gotypes.ToBool(v.Value), v.Usage)
	case ValueTypeInt:
		v.Value = flag.Int(v.Key, gotypes.ToInt(v.Value), v.Usage)
	case ValueTypeInt64:
		v.Value = flag.Int64(v.Key, gotypes.ToInt64(v.Value), v.Usage)
	case ValueTypeUint:
		v.Value = flag.Uint(v.Key, gotypes.ToUint(v.Value), v.Usage)
	case ValueTypeUint64:
		v.Value = flag.Uint64(v.Key, gotypes.ToUint64(v.Value), v.Usage)
	case ValueTypeFloat64:
		v.Value = flag.Float64(v.Key, gotypes.ToFloat64(v.Value), v.Usage)
	case ValueTypeString:
		v.Value = flag.String(v.Key, gotypes.ToString(v.Value), v.Usage)
	case ValueTypeDuration:
		v.Value = flag.Duration(v.Key, gotypes.ToDuration(v.Value), v.Usage)
	}

	c.mutex.Lock()
	c.variables[v.Key] = v
	c.mutex.Unlock()
}

func (c *Component) WatchVariable(key string, watcher Watcher) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if watchers, ok := c.watchers[key]; ok {
		c.watchers[key] = append(watchers, watcher)
	} else {
		c.watchers[key] = []Watcher{watcher}
	}
}

func (c *Component) Has(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if _, ok := c.variables[key]; ok {
		return true
	}

	return false
}

func (c *Component) Get(key string) interface{} {
	c.mutex.RLock()
	v, ok := c.variables[key]
	c.mutex.RUnlock()

	if ok && v.Value != nil {
		switch v.Type {
		case ValueTypeBool:
			return c.GetBool(key)
		case ValueTypeInt:
			return c.GetInt(key)
		case ValueTypeInt64:
			return c.GetInt64(key)
		case ValueTypeUint:
			return c.GetUint(key)
		case ValueTypeUint64:
			return c.GetUint64(key)
		case ValueTypeFloat64:
			return c.GetFloat64(key)
		case ValueTypeString:
			return c.GetString(key)
		case ValueTypeDuration:
			return c.GetDuration(key)
		}
	}

	return nil
}

func (c *Component) Set(key string, value interface{}) error {
	oldValue := c.Get(key)
	c.mutex.Lock()

	variable, ok := c.variables[key]

	if !ok {
		c.mutex.Unlock()
		return errors.New("Config already parsed. Can't and new variable")
	}

	switch variable.Type {
	case ValueTypeBool:
		value = gotypes.ToBool(value)
	case ValueTypeInt:
		value = gotypes.ToInt(value)
	case ValueTypeInt64:
		value = gotypes.ToInt64(value)
	case ValueTypeUint:
		value = gotypes.ToUint(value)
	case ValueTypeUint64:
		value = gotypes.ToUint64(value)
	case ValueTypeFloat64:
		value = gotypes.ToFloat64(value)
	case ValueTypeString:
		value = gotypes.ToString(value)
	case ValueTypeDuration:
		value = gotypes.ToDuration(value)
	default:
		c.mutex.Unlock()
		return fmt.Errorf("Unknown type %s for config %s", variable.Type, variable.Key)
	}

	variable.Value = value
	c.variables[key] = variable

	watchers := []Watcher{}

	if watchersForAll, ok := c.watchers[WatcherForAll]; ok {
		for _, w := range watchersForAll {
			watchers = append(watchers, w)
		}
	}

	if watchersByKey, ok := c.watchers[key]; ok {
		for _, w := range watchersByKey {
			watchers = append(watchers, w)
		}
	}

	c.mutex.Unlock()

	if ok {
		go func() {
			for _, watcher := range watchers {
				watcher(key, value, oldValue)
			}
		}()
	}

	return nil
}

func (c *Component) IsEditable(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if variable, ok := c.variables[key]; ok {
		return variable.Editable
	}

	return false
}

func (c *Component) GetAllVariables() map[string]Variable {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	variables := make(map[string]Variable, len(c.variables))

	for k, v := range c.variables {
		variables[k] = v
	}

	return variables
}

func (c *Component) GetAllValues() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	values := make(map[string]interface{}, len(c.variables))

	for k, v := range c.variables {
		values[k] = v.Value
	}

	return values
}

func (c *Component) GetGlobalConf() *globalconf.GlobalConf {
	return c.config
}

func (c *Component) GetBool(key string) bool {
	return c.GetBoolDefault(key, false)
}

func (c *Component) GetBoolDefault(key string, value interface{}) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if val, ok := c.variables[key]; ok {
		return gotypes.ToBool(val.Value)
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
		return gotypes.ToInt(val.Value)
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
		return gotypes.ToInt64(val.Value)
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
		return gotypes.ToUint(val.Value)
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
		return gotypes.ToUint64(val.Value)
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
		return gotypes.ToFloat64(val.Value)
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
		return gotypes.ToString(val.Value)
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
		return gotypes.ToDuration(val.Value)
	}

	return gotypes.ToDuration(value)
}
