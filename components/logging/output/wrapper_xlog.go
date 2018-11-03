package output

import (
	"path"
	"runtime"
	"strconv"
	"sync"

	"github.com/kihamo/shadow/components/logging"
	"github.com/rs/xlog"
)

type WrapperXLog struct {
	x      xlog.Logger
	config xlog.Config

	mutex sync.RWMutex
}

func ConvertLoggerToXLogLevel(l logging.Level) xlog.Level {
	switch l {
	case logging.LevelEmergency:
		return xlog.LevelFatal
	case logging.LevelAlert:
		return xlog.LevelFatal
	case logging.LevelCritical:
		return xlog.LevelFatal
	case logging.LevelError:
		return xlog.LevelError
	case logging.LevelWarning:
		return xlog.LevelWarn
	case logging.LevelNotice:
		return xlog.LevelInfo
	case logging.LevelInformational:
		return xlog.LevelInfo
	case logging.LevelDebug:
		return xlog.LevelDebug
	}

	return xlog.LevelInfo
}

func NewWrapperXLog(c xlog.Config) *WrapperXLog {
	l := &WrapperXLog{
		x:      xlog.New(c),
		config: c,
	}

	// free memory
	l.config.Fields = nil

	return l
}

func (l *WrapperXLog) SetLevel(lv xlog.Level) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.config.Level = lv
	l.config.Fields = l.x.GetFields()

	l.x = xlog.New(l.config)

	// free memory
	l.config.Fields = nil
}

func (l *WrapperXLog) SetFields(f map[string]interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.config.Fields = f

	l.x = xlog.New(l.config)

	// free memory
	l.config.Fields = nil
}

func (l *WrapperXLog) setCallFile(c int, f map[string]interface{}) {
	if _, ok := f[xlog.KeyFile]; !ok {
		if _, file, line, ok := runtime.Caller(c); ok {
			f[xlog.KeyFile] = path.Base(file) + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
}

func (l *WrapperXLog) getArguments(c int, v ...interface{}) []interface{} {
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
func (l *WrapperXLog) Write(p []byte) (n int, err error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.x.Write(p)
}

func (l *WrapperXLog) SetField(n string, v interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.SetField(n, v)
}

func (l *WrapperXLog) GetFields() xlog.F {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.x.GetFields()
}

func (l *WrapperXLog) Debug(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Debug(l.getArguments(2, v...)...)
}

func (l *WrapperXLog) Debugf(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Debugf(f, l.getArguments(2, v...)...)
}

func (l *WrapperXLog) Info(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Info(l.getArguments(2, v...)...)
}

func (l *WrapperXLog) Infof(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Infof(f, l.getArguments(2, v...)...)
}

func (l *WrapperXLog) Warn(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Warn(l.getArguments(2, v...)...)
}

func (l *WrapperXLog) Warnf(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Warnf(f, l.getArguments(2, v...)...)
}

func (l *WrapperXLog) Error(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Error(l.getArguments(2, v...)...)
}

func (l *WrapperXLog) Errorf(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Errorf(f, l.getArguments(2, v...)...)
}

func (l *WrapperXLog) Fatal(v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Fatal(l.getArguments(2, v...)...)
}

func (l *WrapperXLog) Fatalf(f string, v ...interface{}) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	l.x.Fatalf(f, l.getArguments(2, v...)...)
}

func (l *WrapperXLog) Output(f int, s string) error {
	l.Info(s)
	return nil
}

func (l *WrapperXLog) OutputF(v xlog.Level, c int, m string, f map[string]interface{}) {
	l.setCallFile(c+2, f)

	l.mutex.RLock()
	defer l.mutex.RUnlock()
	l.x.OutputF(v, c, m, f)
}

func (l *WrapperXLog) Printf(f string, v ...interface{}) {
	l.Infof(f, v...)
}

func (l *WrapperXLog) Print(v ...interface{}) {
	l.Info(v...)
}

func (l *WrapperXLog) Println(v ...interface{}) {
	l.Info(v...)
}

func (l *WrapperXLog) Panic(v ...interface{}) {
	l.Fatal(v...)
}

func (l *WrapperXLog) Panicf(f string, v ...interface{}) {
	l.Fatalf(f, v...)
}

func (l *WrapperXLog) Panicln(v ...interface{}) {
	l.Fatal(v...)
}

func (l *WrapperXLog) Log(v ...interface{}) {
	l.Info(v...)
}
