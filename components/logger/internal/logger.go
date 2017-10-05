package internal

import (
	"path"
	"runtime"
	"strconv"
	"sync"

	"github.com/rs/xlog"
)

type loggerWrapper struct {
	x      xlog.Logger
	config xlog.Config

	mutex sync.RWMutex
}

func newLogger(c xlog.Config) *loggerWrapper {
	l := &loggerWrapper{
		x:      xlog.New(c),
		config: c,
	}

	// free memory
	l.config.Fields = nil

	return l
}

func (l *loggerWrapper) setLevel(lv xlog.Level) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.config.Level = lv
	l.config.Fields = l.x.GetFields()

	l.x = xlog.New(l.config)

	// free memory
	l.config.Fields = nil
}

func (l *loggerWrapper) setFields(f map[string]interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.config.Fields = f

	l.x = xlog.New(l.config)

	// free memory
	l.config.Fields = nil
}

func (l *loggerWrapper) setCallFile(c int, f map[string]interface{}) {
	if _, ok := f[xlog.KeyFile]; !ok {
		if _, file, line, ok := runtime.Caller(c); ok {
			f[xlog.KeyFile] = path.Base(file) + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
}

func (l *loggerWrapper) getArguments(c int, v ...interface{}) []interface{} {
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
func (l *loggerWrapper) Write(p []byte) (n int, err error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.x.Write(p)
}

func (l *loggerWrapper) SetField(n string, v interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.SetField(n, v)
}

func (l *loggerWrapper) GetFields() xlog.F {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.x.GetFields()
}

func (l *loggerWrapper) Debug(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Debug(l.getArguments(2, v...)...)
}

func (l *loggerWrapper) Debugf(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Debugf(f, l.getArguments(2, v...)...)
}

func (l *loggerWrapper) Info(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Info(l.getArguments(2, v...)...)
}

func (l *loggerWrapper) Infof(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Infof(f, l.getArguments(2, v...)...)
}

func (l *loggerWrapper) Warn(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Warn(l.getArguments(2, v...)...)
}

func (l *loggerWrapper) Warnf(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Warnf(f, l.getArguments(2, v...)...)
}

func (l *loggerWrapper) Error(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Error(l.getArguments(2, v...)...)
}

func (l *loggerWrapper) Errorf(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Errorf(f, l.getArguments(2, v...)...)
}

func (l *loggerWrapper) Fatal(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Fatal(l.getArguments(2, v...)...)
}

func (l *loggerWrapper) Fatalf(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Fatalf(f, l.getArguments(2, v...)...)
}

func (l *loggerWrapper) Output(f int, s string) error {
	l.Info(s)
	return nil
}

func (l *loggerWrapper) OutputF(v xlog.Level, c int, m string, f map[string]interface{}) {
	l.setCallFile(c+2, f)

	l.mutex.RLock()
	defer l.mutex.RUnlock()
	l.x.OutputF(v, c, m, f)
}

func (l *loggerWrapper) Printf(f string, v ...interface{}) {
	l.Infof(f, v...)
}

func (l *loggerWrapper) Print(v ...interface{}) {
	l.Info(v...)
}

func (l *loggerWrapper) Println(v ...interface{}) {
	l.Info(v...)
}

func (l *loggerWrapper) Panic(v ...interface{}) {
	l.Fatal(v...)
}

func (l *loggerWrapper) Panicf(f string, v ...interface{}) {
	l.Fatalf(f, v...)
}

func (l *loggerWrapper) Panicln(v ...interface{}) {
	l.Fatal(v...)
}

func (l *loggerWrapper) Log(v ...interface{}) {
	l.Info(v...)
}
