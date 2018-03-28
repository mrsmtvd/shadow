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

	Watch(watcher Watcher, source string)
	Watchers(key string) []Watcher

	Variables() []Variable
	Bool(key string) bool
	BoolDefault(key string, value interface{}) bool
	Int(key string) int
	IntDefault(key string, value interface{}) int
	Int64(key string) int64
	Int64Default(key string, value interface{}) int64
	Uint(key string) uint
	UintDefault(key string, value interface{}) uint
	Uint64(key string) uint64
	Uint64Default(key string, value interface{}) uint64
	Float64(key string) float64
	Float64Default(key string, value interface{}) float64
	String(key string) string
	StringDefault(key string, value interface{}) string
	Duration(key string) time.Duration
	DurationDefault(key string, value interface{}) time.Duration
}
