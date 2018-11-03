package logging

import (
	"sync"

	"go.uber.org/zap"
)

type wrapper struct {
	mutex  sync.RWMutex
	logger *zap.SugaredLogger
}

func newWrapper() *wrapper {
	return &wrapper{
		logger: zap.NewNop().Sugar(),
	}
}

func (w *wrapper) SetLogger(l *zap.SugaredLogger) {
	w.mutex.Lock()
	w.logger = l
	w.mutex.Unlock()
}

func (w *wrapper) Named(name string) Logger {
	return &wrapper{
		logger: w.logger.Named(name),
	}
}

func (w *wrapper) Debug(message string, args ...interface{}) {
	w.logger.Debugw(message, args...)
}

func (w *wrapper) Debugf(template string, args ...interface{}) {
	w.logger.Debugf(template, args...)
}

func (w *wrapper) Info(message string, args ...interface{}) {
	w.logger.Infow(message, args...)
}

func (w *wrapper) Infof(template string, args ...interface{}) {
	w.logger.Infof(template, args...)
}

func (w *wrapper) Warn(message string, args ...interface{}) {
	w.logger.Warnw(message, args...)
}

func (w *wrapper) Warnf(template string, args ...interface{}) {
	w.logger.Warnf(template, args...)
}

func (w *wrapper) Error(message string, args ...interface{}) {
	w.logger.Errorw(message, args...)
}

func (w *wrapper) Errorf(template string, args ...interface{}) {
	w.logger.Errorf(template, args...)
}

func (w *wrapper) Panic(message string, args ...interface{}) {
	w.logger.Panicw(message, args...)
}

func (w *wrapper) Panicf(template string, args ...interface{}) {
	w.logger.Panicf(template, args...)
}

func (w *wrapper) Fatal(message string, args ...interface{}) {
	w.logger.Fatalw(message, args...)
}

func (w *wrapper) Fatalf(template string, args ...interface{}) {
	w.logger.Fatalf(template, args...)
}
