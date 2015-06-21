package system

import (
	"github.com/Sirupsen/logrus"
	"github.com/kihamo/shadow/service/frontend"
)

var loggerHook *LoggerHook

func init() {
	loggerHook = NewLoggerHook()
	logrus.AddHook(loggerHook)
}

type LogsHandler struct {
	frontend.AbstractFrontendHandler
}

func (h *LogsHandler) Handle() {
	loggers := loggerHook.GetRecords()
	log := h.Input.URL.Query().Get("log")

	if _, ok := loggers[log]; ok && log != "" {
		reply := make([]map[string]interface{}, len(loggers[log]))

		for i := range loggers[log] {
			reply[i] = map[string]interface{}{
				"time":    loggers[log][i].Time,
				"message": loggers[log][i].Message,
				"level":   loggers[log][i].Level.String(),
			}
		}

		h.SendJSON(reply)
		return
	}

	h.SetTemplate("logs.tpl.html")
	h.View.Context["PageTitle"] = "Log view"
	h.View.Context["Loggers"] = loggers
}
