package output

import (
	"github.com/kihamo/shadow/components/logging"
	"github.com/rs/xlog"
)

func NewConsoleOutput(level logging.Level, fields map[string]interface{}) logging.Logger {
	config := xlog.Config{
		Output: xlog.NewConsoleOutput(),
		Level:  ConvertLoggerToXLogLevel(level),
		Fields: fields,
	}

	return NewWrapperXLog(config)
}
