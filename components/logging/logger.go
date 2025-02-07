package logging

import (
	"github.com/mrsmtvd/shadow/components/logging/internal/wrapper"
)

type Logger wrapper.Logger

var defaultLogger = wrapper.NewNop("default")

func DefaultLogger() Logger {
	return defaultLogger
}

func DefaultLazyLogger(name string) Logger {
	return NewLazyLogger(DefaultLogger(), name)
}
