package config

const (
	WatcherForAll = "*"
)

type Watcher interface {
	Source() string
	Keys() []string
	Callback(key string, new interface{}, old interface{})
}

type HasWatchers interface {
	GetConfigWatchers() []Watcher
}

type WatcherItem struct {
	source   string
	keys     []string
	callback func(string, interface{}, interface{})
}

func NewWatcher(source string, keys []string, callback func(string, interface{}, interface{})) Watcher {
	return &WatcherItem{
		source:   source,
		keys:     keys,
		callback: callback,
	}
}

func (w *WatcherItem) Source() string {
	return w.source
}

func (w *WatcherItem) Keys() []string {
	return w.keys
}

func (w *WatcherItem) Callback(key string, new interface{}, old interface{}) {
	w.callback(key, new, old)
}
