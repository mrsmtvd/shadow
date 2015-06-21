package resource

import (
	"flag"

	"github.com/kihamo/shadow"
	"github.com/vharitonsky/iniflags"
)

type Config struct {
	application *shadow.Application
	values      map[string]interface{}
}

type ConfigVariable struct {
	Key   string
	Value interface{}
	Usage string
}

type ServiceConfigurable interface {
	GetConfigVariables() []ConfigVariable
}

func (r *Config) GetName() string {
	return "config"
}

func (r *Config) Init(a *shadow.Application) error {
	r.application = a

	r.values = map[string]interface{}{}
	r.Add("debug", false, "Debug mode")
	r.Add("env", "stable", "Environment")

	return nil
}

func (r *Config) Run() error {
	for _, service := range r.application.GetServices() {
		if configurable, ok := service.(ServiceConfigurable); ok {
			for _, variable := range configurable.GetConfigVariables() {
				r.Add(variable.Key, variable.Value, variable.Usage)
			}
		}
	}

	iniflags.Parse()

	resourceLogger, err := r.application.GetResource("logger")
	if err == nil {
		flag.VisitAll(func(f *flag.Flag) {
			if f.Name == "config" {
				resourceLogger.(*Logger).Get(r.GetName()).Infof("Use config from %s", f.Value.String())
			}
		})


	}

	return nil
}

func (r *Config) Add(key string, value interface{}, usage string) {
	switch value.(type) {
	case bool:
		r.values[key] = flag.Bool(key, value.(bool), usage)
	case int64:
		r.values[key] = flag.Int64(key, value.(int64), usage)
	case string:
		r.values[key] = flag.String(key, value.(string), usage)
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
		switch val.(type) {
		case *bool:
			return *(val.(*bool))
		case *int64:
			return *(val.(*int64))
		case *string:
			return *(val.(*string))
		}
	}

	return nil
}

func (r *Config) GetAll() map[string]interface{} {
	return r.values
}

func (r *Config) GetBool(key string) bool {
	if val, ok := r.values[key]; ok {
		if res, okCast := val.(*bool); okCast {
			return *res
		}
	}

	return false
}

func (r *Config) GetInt64(key string) int64 {
	if val, ok := r.values[key]; ok {
		if res, okCast := val.(*int64); okCast {
			return *res
		}
	}

	return -1
}

func (r *Config) GetString(key string) string {
	if val, ok := r.values[key]; ok {
		if res, okCast := val.(*string); okCast {
			return *res
		}
	}

	return ""
}
