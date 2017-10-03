package config

const (
	WatcherForAll = "*"
)

type Watcher func(string, interface{}, interface{})

type HasWatchers interface {
	GetConfigWatchers() map[string][]Watcher
}
