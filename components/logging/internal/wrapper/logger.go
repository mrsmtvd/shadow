package wrapper

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
}
