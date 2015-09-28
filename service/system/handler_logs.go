package system

import (
	"github.com/Sirupsen/logrus"
	"github.com/kihamo/shadow/service/frontend"
)

var LoggerHookInstance *LoggerHook

func init() {
	LoggerHookInstance = NewLoggerHook()
	logrus.AddHook(LoggerHookInstance)
}

type LogsHandler struct {
	frontend.AbstractFrontendHandler
}

func (h *LogsHandler) Handle() {
	loggers := LoggerHookInstance.GetRecords()
	log := h.Input.URL.Query().Get("log")

	if _, ok := loggers[log]; ok && log != "" {
		reply := make([]map[string]interface{}, len(loggers[log]))

		for i, logger := range loggers[log] {
			fields := map[string]interface{}{}
			for name, value := range logger.Data {
				if name != "component" {
					fields[name] = value
				}
			}

			reply[i] = map[string]interface{}{
				"time":    logger.Time,
				"message": logger.Message,
				"level":   logger.Level.String(),
				"fields":  fields,
			}
		}

		h.SendJSON(reply)
		return
	}

	h.SetTemplate("logs.tpl.html")
	h.SetPageTitle("Log view")
	h.SetPageHeader("Log view")
	h.SetVar("Loggers", loggers)
}
