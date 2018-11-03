package internal

import (
	"github.com/kihamo/shadow/components/logging"
)

type logger struct {
	l logging.Logger
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.l.Infof(format, v...)
}
