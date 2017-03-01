package logger

import (
	"github.com/kihamo/shadow"
)

func NewOrNop(name string, application shadow.Application) Logger {
	if resourceLogger, err := application.GetComponent(ComponentName); err == nil {
		return resourceLogger.(*Component).Get(name)
	}

	return NopLogger
}
