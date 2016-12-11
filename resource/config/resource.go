package config

import (
	"errors"
	"flag"
	"strings"
	"sync"
	"time"

	"github.com/kihamo/gotypes"
	"github.com/kihamo/shadow"
	"github.com/rakyll/globalconf"
)

const (
	FlagConfig = "config"

	ValueTypeBool     = "bool"
	ValueTypeInt      = "int"
	ValueTypeInt64    = "int64"
	ValueTypeUint     = "uint"
	ValueTypeUint64   = "uint64"
	ValueTypeFloat64  = "float64"
	ValueTypeString   = "string"
	ValueTypeDuration = "duration"
)

type Watcher func(interface{}, interface{})

type Resource struct {
	mutex       sync.RWMutex
	application *shadow.Application
	config      *globalconf.GlobalConf
	variables   map[string]Variable
	watchers    map[string][]Watcher
}

type Variable struct {
	Key      string
	Default  interface{}
	Value    interface{}
	Type     string
	Usage    string
	Editable bool
}

type hasVariables interface {
	GetConfigVariables() []Variable
}

type hasWatchers interface {
	GetConfigWatcher() map[string][]Watcher
}

func (r *Resource) GetName() string {
	return "config"
}

func (r *Resource) Init(a *shadow.Application) (err error) {
	r.application = a

	config := flag.String(FlagConfig, "", "Config file which which override default config parameters")
	flag.Parse()

	opts := globalconf.Options{
		EnvPrefix: strings.ToUpper(strings.Replace(a.Name, " ", "_", -1)) + "_",
		Filename:  *config,
	}

	if r.config, err = globalconf.NewWithOptions(&opts); err != nil {
		return err
	}

	r.mutex.Lock()
	r.variables = map[string]Variable{}
	r.watchers = map[string][]Watcher{}
	r.mutex.Unlock()

	return err
}

func (r *Resource) Run() error {
	for _, resource := range r.application.GetResources() {
		if variables, ok := resource.(hasVariables); ok {
			for _, variable := range variables.GetConfigVariables() {
				r.addFlag(variable)
			}
		}

		if watchers, ok := resource.(hasWatchers); ok {
			for key, list := range watchers.GetConfigWatcher() {
				for _, watcher := range list {
					r.WatchVariable(key, watcher)
				}
			}
		}
	}

	for _, service := range r.application.GetServices() {
		if variables, ok := service.(hasVariables); ok {
			for _, variable := range variables.GetConfigVariables() {
				r.addFlag(variable)
			}
		}

		if watchers, ok := service.(hasWatchers); ok {
			for key, list := range watchers.GetConfigWatcher() {
				for _, watcher := range list {
					r.WatchVariable(key, watcher)
				}
			}
		}
	}

	r.config.ParseAll()

	return nil
}

func (r *Resource) addFlag(v Variable) {
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

	r.mutex.Lock()
	r.variables[v.Key] = v
	r.mutex.Unlock()
}

func (r *Resource) WatchVariable(key string, watcher Watcher) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if watchers, ok := r.watchers[key]; ok {
		r.watchers[key] = append(watchers, watcher)
	} else {
		r.watchers[key] = []Watcher{watcher}
	}
}

func (r *Resource) Has(key string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if _, ok := r.variables[key]; ok {
		return true
	}

	return false
}

func (r *Resource) Get(key string) interface{} {
	r.mutex.RLock()
	v, ok := r.variables[key]
	r.mutex.RUnlock()

	if ok && v.Value != nil {
		switch v.Type {
		case ValueTypeBool:
			return r.GetBool(key)
		case ValueTypeInt:
			return r.GetInt(key)
		case ValueTypeInt64:
			return r.GetInt64(key)
		case ValueTypeUint:
			return r.GetUint(key)
		case ValueTypeUint64:
			return r.GetUint64(key)
		case ValueTypeFloat64:
			return r.GetFloat64(key)
		case ValueTypeString:
			return r.GetString(key)
		case ValueTypeDuration:
			return r.GetDuration(key)
		}
	}

	return nil
}

func (r *Resource) Set(key string, value interface{}) error {
	old := r.Get(key)
	r.mutex.Lock()

	variable, ok := r.variables[key]

	if !ok {
		return errors.New("Config already parsed. Can't and new variable")
	}

	variable.Value = value

	r.variables[key] = variable
	watchers, ok := r.watchers[key]

	r.mutex.Unlock()

	if ok {
		go func() {
			for _, watcher := range watchers {
				watcher(value, old)
			}
		}()
	}

	return nil
}

func (r *Resource) IsEditable(key string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if variable, ok := r.variables[key]; ok {
		return variable.Editable
	}

	return false
}

func (r *Resource) GetAllVariables() map[string]Variable {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	variables := make(map[string]Variable, len(r.variables))

	for k, v := range r.variables {
		variables[k] = v
	}

	return variables
}

func (r *Resource) GetAllValues() map[string]interface{} {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	values := make(map[string]interface{}, len(r.variables))

	for k, v := range r.variables {
		values[k] = v.Value
	}

	return values
}

func (r *Resource) GetGlobalConf() *globalconf.GlobalConf {
	return r.config
}

func (r *Resource) GetBool(key string) bool {
	return r.GetBoolDefault(key, false)
}

func (r *Resource) GetBoolDefault(key string, value interface{}) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if val, ok := r.variables[key]; ok {
		return gotypes.ToBool(val)
	}

	return gotypes.ToBool(value)
}

func (r *Resource) GetInt(key string) int {
	return r.GetIntDefault(key, -1)
}

func (r *Resource) GetIntDefault(key string, value interface{}) int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if val, ok := r.variables[key]; ok {
		return gotypes.ToInt(val.Value)
	}

	return gotypes.ToInt(value)
}

func (r *Resource) GetInt64(key string) int64 {
	return r.GetInt64Default(key, -1)
}

func (r *Resource) GetInt64Default(key string, value interface{}) int64 {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if val, ok := r.variables[key]; ok {
		return gotypes.ToInt64(val.Value)
	}

	return gotypes.ToInt64(value)
}

func (r *Resource) GetUint(key string) uint {
	return r.GetUintDefault(key, 0)
}

func (r *Resource) GetUintDefault(key string, value interface{}) uint {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if val, ok := r.variables[key]; ok {
		return gotypes.ToUint(val.Value)
	}

	return gotypes.ToUint(value)
}

func (r *Resource) GetUint64(key string) uint64 {
	return r.GetUint64Default(key, 0)
}

func (r *Resource) GetUint64Default(key string, value interface{}) uint64 {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if val, ok := r.variables[key]; ok {
		return gotypes.ToUint64(val.Value)
	}

	return gotypes.ToUint64(value)
}

func (r *Resource) GetFloat64(key string) float64 {
	return r.GetFloat64Default(key, -1)
}

func (r *Resource) GetFloat64Default(key string, value interface{}) float64 {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if val, ok := r.variables[key]; ok {
		return gotypes.ToFloat64(val.Value)
	}

	return gotypes.ToFloat64(value)
}

func (r *Resource) GetString(key string) string {
	return r.GetStringDefault(key, "")
}

func (r *Resource) GetStringDefault(key string, value interface{}) string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if val, ok := r.variables[key]; ok {
		return gotypes.ToString(val.Value)
	}

	return gotypes.ToString(value)
}

func (r *Resource) GetDuration(key string) time.Duration {
	return r.GetDurationDefault(key, 0)
}

func (r *Resource) GetDurationDefault(key string, value interface{}) time.Duration {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if val, ok := r.variables[key]; ok {
		return gotypes.ToDuration(val.Value)
	}

	return gotypes.ToDuration(value)
}
