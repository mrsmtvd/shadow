package resource

import (
	"flag"
	"reflect"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kihamo/gotypes"
	"github.com/kihamo/shadow"
	"github.com/rakyll/globalconf"
)

const (
	flagConfig = "config"
)

type Config struct {
	application *shadow.Application
	config      *globalconf.GlobalConf
	values      map[string]interface{}
}

type ConfigVariable struct {
	Key   string
	Value interface{}
	Usage string
}

type ContextItemConfigurable interface {
	GetConfigVariables() []ConfigVariable
}

func (r *Config) GetName() string {
	return "config"
}

func (r *Config) Init(a *shadow.Application) (err error) {
	r.application = a

	config := flag.String(flagConfig, "", "Config file which which override default config parameters")
	flag.Parse()

	opts := globalconf.Options{
		EnvPrefix: strings.ToUpper(strings.Replace(a.Name, " ", "_", -1)) + "_",
		Filename:  *config,
	}

	if r.config, err = globalconf.NewWithOptions(&opts); err != nil {
		return err
	}

	r.values = map[string]interface{}{}

	return err
}

func (r *Config) Run() error {
	r.Add("debug", false, "Debug mode")

	for _, resource := range r.application.GetResources() {
		if configurable, ok := resource.(ContextItemConfigurable); ok {
			for _, variable := range configurable.GetConfigVariables() {
				r.Add(variable.Key, variable.Value, variable.Usage)
			}
		}
	}

	for _, service := range r.application.GetServices() {
		if configurable, ok := service.(ContextItemConfigurable); ok {
			for _, variable := range configurable.GetConfigVariables() {
				r.Add(variable.Key, variable.Value, variable.Usage)
			}
		}
	}

	r.config.ParseAll()

	resourceLogger, err := r.application.GetResource("logger")
	if err == nil {
		fields := logrus.Fields{}
		for key := range r.GetAll() {
			fields[key] = r.Get(key)
		}

		logger := resourceLogger.(*Logger).Get(r.GetName())
		logger.WithFields(fields).Infof("Config env prefix %s", r.config.EnvPrefix)

		flag.VisitAll(func(f *flag.Flag) {
			if f.Name == flagConfig && f.Value.String() != "" {
				logger.Infof("Use config from %s", f.Value.String())
			}
		})
	}

	return nil
}

func (r *Config) Add(key string, value interface{}, usage string) {
	switch value.(type) {
	case bool:
		r.values[key] = flag.Bool(key, value.(bool), usage)
	case int:
		r.values[key] = flag.Int(key, value.(int), usage)
	case int64:
		r.values[key] = flag.Int64(key, value.(int64), usage)
	case uint:
		r.values[key] = flag.Uint(key, value.(uint), usage)
	case uint64:
		r.values[key] = flag.Uint64(key, value.(uint64), usage)
	case float64:
		r.values[key] = flag.Float64(key, value.(float64), usage)
	case string:
		r.values[key] = flag.String(key, value.(string), usage)
	case time.Duration:
		r.values[key] = flag.Duration(key, value.(time.Duration), usage)
	}
}

func (r *Config) Has(key string) bool {
	if _, ok := r.values[key]; ok {
		return true
	}

	return false
}

func (r *Config) Get(key string) interface{} {
	if val, ok := r.values[key]; ok {
		return reflect.Indirect(reflect.ValueOf(val)).Interface()
	}

	return nil
}

func (r *Config) GetAll() map[string]interface{} {
	return r.values
}

func (r *Config) GetBool(key string) bool {
	if val, ok := r.values[key]; ok {
		return gotypes.ToBool(val)
	}

	return false
}

func (r *Config) GetInt(key string) int {
	if val, ok := r.values[key]; ok {
		return gotypes.ToInt(val)
	}

	return -1
}

func (r *Config) GetInt64(key string) int64 {
	if val, ok := r.values[key]; ok {
		return gotypes.ToInt64(val)
	}

	return -1
}

func (r *Config) GetUint(key string) uint {
	if val, ok := r.values[key]; ok {
		return gotypes.ToUint(val)
	}

	return 0
}

func (r *Config) GetUint64(key string) uint64 {
	if val, ok := r.values[key]; ok {
		return gotypes.ToUint64(val)
	}

	return 0
}

func (r *Config) GetFloat64(key string) float64 {
	if val, ok := r.values[key]; ok {
		return gotypes.ToFloat64(val)
	}

	return -1
}

func (r *Config) GetString(key string) string {
	if val, ok := r.values[key]; ok {
		return gotypes.ToString(val)
	}

	return ""
}

func (r *Config) GetDuration(key string) time.Duration {
	if val := r.GetString(key); val != "" {
		if r, err := time.ParseDuration(val); err == nil {
			return r
		}
	}

	return 0
}
