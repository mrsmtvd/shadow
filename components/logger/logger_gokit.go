package logger

import (
	"github.com/go-kit/kit/log"
)

type GoKitLogger struct {
	l Logger
}

func NewGoKitLogger(l Logger) log.Logger {
	return GoKitLogger{
		l: l,
	}
}

func (l GoKitLogger) Log(v ...interface{}) error {
	l.l.Log(v)
	return nil
}
