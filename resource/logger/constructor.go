package logger

import (
	"github.com/kihamo/shadow"
)

func NewOrNop(name string, application *shadow.Application) Logger {
	if resourceLogger, err := application.GetResource("logger"); err == nil {
		return resourceLogger.(*Resource).Get(name)
	}

	return NopLogger
}
