package logger

import (
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	Get(key string) Logger
}
