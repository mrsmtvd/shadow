package database

import (
	"github.com/go-gorp/gorp"
	"github.com/rs/xlog"
)

type databaseLogger struct {
	logger xlog.Logger
}

func newDatabaseLogger(l xlog.Logger) gorp.GorpLogger {
	return databaseLogger{
		logger: l,
	}
}

func (l databaseLogger) Printf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
}
