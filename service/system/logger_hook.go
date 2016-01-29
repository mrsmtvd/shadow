package system

import (
	"errors"
	"sync"

	"github.com/Sirupsen/logrus"
)

const (
	MaxItems = 100
)

type LoggerHook struct {
	records map[string][]*logrus.Entry
	mu      sync.Mutex
}

func NewLoggerHook() *LoggerHook {
	return &LoggerHook{
		records: map[string][]*logrus.Entry{},
	}
}

func (h *LoggerHook) GetRecords() map[string][]*logrus.Entry {
	return h.records
}

func (h *LoggerHook) Fire(e *logrus.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, ok := e.Data["component"]
	if !ok {
		return errors.New("Component field not found in log entry")
	}

	component := e.Data["component"].(string)
	if _, ok := h.records[component]; !ok {
		h.records[component] = make([]*logrus.Entry, 0, MaxItems)
	}

	entry := *e
	if len(h.records[component]) == cap(h.records[component]) {
		h.records[component] = append(h.records[component][1:], []*logrus.Entry{&entry}...)
	} else {
		h.records[component] = append(h.records[component], []*logrus.Entry{&entry}...)
	}

	return nil
}

func (h *LoggerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
