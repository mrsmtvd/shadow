package logging

import (
	"github.com/kihamo/shadow/components/logging/internal/wrapper"
)

type Logger wrapper.Logger

var defaultLogger = wrapper.NewNop("default")

func DefaultLogger() Logger {
	return defaultLogger
}
