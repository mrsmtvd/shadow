package config

const (
	WatcherForAll = "*"
)

type Watcher interface {
	Keys() []string
	Callback(key string, new interface{}, old interface{})
}

type HasWatchers interface {
	ConfigWatchers() []Watcher
}

type WatcherSimple struct {
	keys     []string
	callback func(string, interface{}, interface{})
}

func NewWatcher(keys []string, callback func(string, interface{}, interface{})) *WatcherSimple {
	return &WatcherSimple{
		keys:     keys,
		callback: callback,
	}
}

func (w *WatcherSimple) Keys() []string {
	return w.keys
}

func (w *WatcherSimple) Callback(key string, new interface{}, old interface{}) {
	w.callback(key, new, old)
}
