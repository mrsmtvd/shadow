package logger

import (
	"github.com/rs/xlog"
)

type Logger interface {
	xlog.Logger

	Printf(string, ...interface{})
	Print(...interface{})
	Println(...interface{})
	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
	Log(...interface{}) error
}

type logger struct {
	x xlog.Logger
}

func (l *logger) Write(p []byte) (n int, err error) {
	return l.x.Write(p)
}

func (l *logger) SetField(n string, v interface{}) {
	l.x.SetField(n, v)
}

func (l *logger) Debug(v ...interface{}) {
	l.x.Debug(v...)
}

func (l *logger) Debugf(f string, v ...interface{}) {
	l.x.Debugf(f, v...)
}

func (l *logger) Info(v ...interface{}) {
	l.x.Info(v...)
}

func (l *logger) Infof(f string, v ...interface{}) {
	l.x.Infof(f, v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.x.Warn(v...)
}

func (l *logger) Warnf(f string, v ...interface{}) {
	l.x.Warnf(f, v...)
}

func (l *logger) Error(v ...interface{}) {
	l.x.Error(v...)
}

func (l *logger) Errorf(f string, v ...interface{}) {
	l.x.Errorf(f, v...)
}

func (l *logger) Fatal(v ...interface{}) {
	l.x.Fatal(v...)
}

func (l *logger) Fatalf(f string, v ...interface{}) {
	l.x.Fatalf(f, v...)
}

func (l *logger) Output(f int, s string) error {
	return l.x.Output(f, s)
}

func (l *logger) OutputF(v xlog.Level, c int, m string, f map[string]interface{}) {
	l.x.OutputF(v, c, m, f)
}

func (l *logger) Printf(f string, v ...interface{}) {
	l.Infof(f, v...)
}

func (l *logger) Print(v ...interface{}) {
	l.Info(v...)
}

func (l *logger) Println(v ...interface{}) {
	l.Info(v...)
}

func (l *logger) Panic(v ...interface{}) {
	l.Fatal(v...)
}

func (l *logger) Panicf(f string, v ...interface{}) {
	l.Fatalf(f, v...)
}

func (l *logger) Panicln(v ...interface{}) {
	l.x.Fatal(v...)
}

func (l *logger) Log(v ...interface{}) error {
	l.x.Info(v...)
	return nil
}
