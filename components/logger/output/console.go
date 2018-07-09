package output

import (
	"github.com/kihamo/shadow/components/logger"
	"github.com/rs/xlog"
)

func NewConsoleOutput(level logger.Level, fields map[string]interface{}) logger.Logger {
	config := xlog.Config{
		Output: xlog.NewConsoleOutput(),
		Level:  ConvertLoggerToXLogLevel(level),
		Fields: fields,
	}

	return NewWrapperXLog(config)
}
