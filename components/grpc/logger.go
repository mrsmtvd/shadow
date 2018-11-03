package grpc

import (
	"github.com/kihamo/shadow/components/logging"
	"google.golang.org/grpc/grpclog"
)

type Logger struct {
	grpclog.LoggerV2

	logger logging.Logger
}

func NewLogger(l logging.Logger) *Logger {
	return &Logger{
		logger: l,
	}
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *Logger) Warning(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *Logger) Warningln(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *Logger) V(level int) bool {
	return true
}
