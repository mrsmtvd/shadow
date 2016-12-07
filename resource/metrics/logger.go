package metrics

import (
	"github.com/go-kit/kit/log"
	"github.com/rs/xlog"
)

type metricsLogger struct {
	logger xlog.Logger
}

func newMetricsLogger(l xlog.Logger) log.Logger {
	return metricsLogger{
		logger: l,
	}
}

func (l metricsLogger) Log(keyvals ...interface{}) error {
	l.logger.Error(keyvals...)

	return nil
}
