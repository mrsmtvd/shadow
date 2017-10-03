package config

import (
	"time"

	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	EnvPrefix() string

	Has(key string) bool
	Get(key string) interface{}
	Set(key string, value interface{}) error
	IsEditable(key string) bool

	Watch(watcher Watcher)
	GetWatchers(key string) []Watcher

	GetAllVariables() map[string]Variable
	GetBool(key string) bool
	GetBoolDefault(key string, value interface{}) bool
	GetInt(key string) int
	GetIntDefault(key string, value interface{}) int
	GetInt64(key string) int64
	GetInt64Default(key string, value interface{}) int64
	GetUint(key string) uint
	GetUintDefault(key string, value interface{}) uint
	GetUint64Default(key string, value interface{}) uint64
	GetFloat64(key string) float64
	GetFloat64Default(key string, value interface{}) float64
	GetString(key string) string
	GetStringDefault(key string, value interface{}) string
	GetDuration(key string) time.Duration
	GetDurationDefault(key string, value interface{}) time.Duration
}
