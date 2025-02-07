package internal

import (
	"github.com/mrsmtvd/shadow/components/logging"
)

type logger struct {
	l logging.Logger
}

func (l *logger) Println(v ...interface{}) {
	l.l.Info("", v...)
}
