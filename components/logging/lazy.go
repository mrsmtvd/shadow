package logging

import (
	"sync"

	"github.com/kihamo/shadow/components/logging/internal/wrapper"
	"go.uber.org/zap"
)

type LazyLogger struct {
	parent   Logger
	name     string
	once     sync.Once
	instance Logger
}

func NewLazyLogger(parent Logger, name string) Logger {
	if parent == nil {
		parent = DefaultLogger()
	}

	return &LazyLogger{
		parent: parent,
		name:   name,
	}
}

func (l *LazyLogger) logger() Logger {
	l.once.Do(func() {
		l.instance = l.parent.Named(l.name)

		if w, ok := l.instance.(*wrapper.Wrapper); ok {
			w.WithOptions(false, zap.AddCallerSkip(1))
		}
	})

	return l.instance
}

func (l *LazyLogger) Name() string {
	return l.name
}

func (l *LazyLogger) Named(name string) wrapper.Logger {
	return l.logger().Named(name)
}

func (l *LazyLogger) Debug(message string, args ...interface{}) {
	l.logger().Debug(message, args...)
}

func (l *LazyLogger) Debugf(message string, args ...interface{}) {
	l.logger().Debugf(message, args...)
}

func (l *LazyLogger) Info(message string, args ...interface{}) {
	l.logger().Info(message, args...)
}

func (l *LazyLogger) Infof(message string, args ...interface{}) {
	l.logger().Infof(message, args...)
}

func (l *LazyLogger) Warn(message string, args ...interface{}) {
	l.logger().Warn(message, args...)
}

func (l *LazyLogger) Warnf(message string, args ...interface{}) {
	l.logger().Warnf(message, args...)
}

func (l *LazyLogger) Error(message string, args ...interface{}) {
	l.logger().Error(message, args...)
}

func (l *LazyLogger) Errorf(message string, args ...interface{}) {
	l.logger().Errorf(message, args...)
}

func (l *LazyLogger) Panic(message string, args ...interface{}) {
	l.logger().Panic(message, args...)
}

func (l *LazyLogger) Panicf(message string, args ...interface{}) {
	l.logger().Panicf(message, args...)
}

func (l *LazyLogger) Fatal(message string, args ...interface{}) {
	l.logger().Fatal(message, args...)
}

func (l *LazyLogger) Fatalf(message string, args ...interface{}) {
	l.logger().Fatalf(message, args...)
}
