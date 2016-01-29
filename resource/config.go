package resource

import (
	"flag"
	"strings"

	"github.com/Sirupsen/logrus"
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
		r.values[key] = flag.Int64(key, int64(value.(int)), usage)
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
