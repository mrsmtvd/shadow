package logger

import (
	"github.com/kihamo/shadow"
)

func NewOrNop(name string, application shadow.Application) Logger {
	if cmp := application.GetComponent(ComponentName); cmp != nil {
		return cmp.(*Component).Get(name)
	}

	return NopLogger
}
