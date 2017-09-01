package logger

import (
	"path"
	"runtime"
	"strconv"
	"sync"

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
	Log(...interface{})
}

type logger struct {
	x      xlog.Logger
	config xlog.Config

	mutex sync.RWMutex
}

func newLogger(c xlog.Config) *logger {
	l := &logger{
		x:      xlog.New(c),
		config: c,
	}

	// free memory
	l.config.Fields = nil

	return l
}

func (l *logger) setLevel(lv xlog.Level) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.config.Level = lv
	l.config.Fields = l.x.GetFields()

	l.x = xlog.New(l.config)

	// free memory
	l.config.Fields = nil
}

func (l *logger) setFields(f map[string]interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.config.Fields = f

	l.x = xlog.New(l.config)

	// free memory
	l.config.Fields = nil
}

func (l *logger) setCallFile(c int, f map[string]interface{}) {
	if _, ok := f[xlog.KeyFile]; !ok {
		if _, file, line, ok := runtime.Caller(c); ok {
			f[xlog.KeyFile] = path.Base(file) + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
}

func (l *logger) getArguments(c int, v ...interface{}) []interface{} {
	fields := map[string]interface{}{}
	args := make([]interface{}, 0)

	// merge field maps
	for i := range v {
		if set, ok := v[i].(map[string]interface{}); ok {
			for key, val := range set {
				fields[key] = val
			}
		} else {
			args = append(args, v[i])
		}
	}

	l.setCallFile(c+1, fields)

	return append(args, fields)
}

// FIXME: file field
func (l *logger) Write(p []byte) (n int, err error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.x.Write(p)
}

func (l *logger) SetField(n string, v interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.SetField(n, v)
}

func (l *logger) GetFields() xlog.F {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.x.GetFields()
}

func (l *logger) Debug(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Debug(l.getArguments(2, v...)...)
}

func (l *logger) Debugf(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Debugf(f, l.getArguments(2, v...)...)
}

func (l *logger) Info(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Info(l.getArguments(2, v...)...)
}

func (l *logger) Infof(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Infof(f, l.getArguments(2, v...)...)
}

func (l *logger) Warn(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Warn(l.getArguments(2, v...)...)
}

func (l *logger) Warnf(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Warnf(f, l.getArguments(2, v...)...)
}

func (l *logger) Error(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Error(l.getArguments(2, v...)...)
}

func (l *logger) Errorf(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Errorf(f, l.getArguments(2, v...)...)
}

func (l *logger) Fatal(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Fatal(l.getArguments(2, v...)...)
}

func (l *logger) Fatalf(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Fatalf(f, l.getArguments(2, v...)...)
}

func (l *logger) Output(f int, s string) error {
	l.Info(s)
	return nil
}

func (l *logger) OutputF(v xlog.Level, c int, m string, f map[string]interface{}) {
	l.setCallFile(c+2, f)

	l.mutex.RLock()
	defer l.mutex.RUnlock()
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
	l.Fatal(v...)
}

func (l *logger) Log(v ...interface{}) {
	l.Info(v...)
}
