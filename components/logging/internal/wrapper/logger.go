package wrapper

import (
	"go.uber.org/zap/zapcore"
)

type LevelEnabler zapcore.LevelEnabler

const (
	DebugLevel  = zapcore.DebugLevel
	InfoLevel   = zapcore.InfoLevel
	WarnLevel   = zapcore.WarnLevel
	ErrorLevel  = zapcore.ErrorLevel
	DPanicLevel = zapcore.DPanicLevel
	PanicLevel  = zapcore.PanicLevel
	FatalLevel  = zapcore.FatalLevel
)

type Logger interface {
	Name() string
	Named(string) Logger
	Debug(string, ...interface{})
	Debugf(string, ...interface{})
	Info(string, ...interface{})
	Infof(string, ...interface{})
	Warn(string, ...interface{})
	Warnf(string, ...interface{})
	Error(string, ...interface{})
	Errorf(string, ...interface{})
	Panic(string, ...interface{})
	Panicf(string, ...interface{})
	Fatal(string, ...interface{})
	Fatalf(string, ...interface{})
	SetLevelEnabler(full bool, level LevelEnabler)
}
