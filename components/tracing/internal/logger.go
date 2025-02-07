package internal

import (
	"strings"

	"github.com/mrsmtvd/shadow/components/logging"
	"github.com/uber/jaeger-client-go"
)

type Logger struct {
	jaeger.Logger
	logger logging.Logger
}

func NewLogger(l logging.Logger) *Logger {
	return &Logger{
		logger: l,
	}
}

func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.logger.Infof(strings.TrimSpace(msg), args...)
}
