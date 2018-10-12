package internal

import (
	"strings"

	"github.com/kihamo/shadow/components/logger"
	"github.com/uber/jaeger-client-go"
)

type Logger struct {
	jaeger.Logger
	logger logger.Logger
}

func NewLogger(l logger.Logger) *Logger {
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
