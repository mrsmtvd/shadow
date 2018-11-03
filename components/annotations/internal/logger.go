package internal

import (
	"github.com/kihamo/shadow/components/logging"
)

type logger struct {
	l logging.Logger
}

func (l *logger) Println(v ...interface{}) {
	l.l.Info("", v...)
}
