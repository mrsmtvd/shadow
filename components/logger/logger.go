package logger

import (
	"github.com/rs/xlog"
)

type Level int

const (
	LevelEmergency Level = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

type Logger interface {
	xlog.Logger

	Printf(string, ...interface{})
	Print(...interface{})
	Println(...interface{})
	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
	Log(...interface{})
}
