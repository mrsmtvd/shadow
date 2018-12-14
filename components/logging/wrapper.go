package logging

import (
	"sync"

	"go.uber.org/zap"
)

type wrapper struct {
	mutex  sync.RWMutex
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	sub    []*wrapper
	name   string
}

func newWrapper() *wrapper {
	return newWrapperByLogger(zap.L())
}

func newWrapperByLogger(l *zap.Logger) *wrapper {
	w := &wrapper{
		sub: make([]*wrapper, 0),
	}
	w.SetLogger(l)

	return w
}

func (w *wrapper) Sugar() *zap.SugaredLogger {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return w.sugar
}

func (w *wrapper) Logger() *zap.Logger {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return w.logger
}

func (w *wrapper) SetLogger(l *zap.Logger) {
	w.mutex.Lock()

	w.logger = l
	w.sugar = l.Sugar()

	for _, sub := range w.sub {
		sub.SetLogger(l.Named(sub.name))
	}

	w.mutex.Unlock()
}

func (w *wrapper) Named(name string) Logger {
	newLogger := newWrapperByLogger(w.Logger().Named(name))
	newLogger.name = name

	w.mutex.Lock()
	w.sub = append(w.sub, newLogger)
	w.mutex.Unlock()

	return newLogger
}

func (w *wrapper) Debug(message string, args ...interface{}) {
	w.Sugar().Debugw(message, args...)
}

func (w *wrapper) Debugf(template string, args ...interface{}) {
	w.Sugar().Debugf(template, args...)
}

func (w *wrapper) Info(message string, args ...interface{}) {
	w.Sugar().Infow(message, args...)
}

func (w *wrapper) Infof(template string, args ...interface{}) {
	w.Sugar().Infof(template, args...)
}

func (w *wrapper) Warn(message string, args ...interface{}) {
	w.Sugar().Warnw(message, args...)
}

func (w *wrapper) Warnf(template string, args ...interface{}) {
	w.Sugar().Warnf(template, args...)
}

func (w *wrapper) Error(message string, args ...interface{}) {
	w.Sugar().Errorw(message, args...)
}

func (w *wrapper) Errorf(template string, args ...interface{}) {
	w.Sugar().Errorf(template, args...)
}

func (w *wrapper) Panic(message string, args ...interface{}) {
	w.Sugar().Panicw(message, args...)
}

func (w *wrapper) Panicf(template string, args ...interface{}) {
	w.Sugar().Panicf(template, args...)
}

func (w *wrapper) Fatal(message string, args ...interface{}) {
	w.Sugar().Fatalw(message, args...)
}

func (w *wrapper) Fatalf(template string, args ...interface{}) {
	w.Sugar().Fatalf(template, args...)
}
