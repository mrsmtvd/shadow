package dashboard

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/kihamo/shadow/components/logger"
)

type PanicHandler struct {
	TemplateHandler

	logger logger.Logger
	error  interface{}
}

func (h *PanicHandler) SetError(err interface{}) {
	h.error = err
}

func (h *PanicHandler) Handle() {
	h.SetView("dashboard", "500")

	h.Response().WriteHeader(http.StatusInternalServerError)

	stack := make([]byte, 4096)
	stack = stack[:runtime.Stack(stack, false)]

	_, filePath, line, _ := runtime.Caller(0)

	fields := map[string]interface{}{
		"error": fmt.Sprintf("%s", h.error),
		"stack": string(stack),
		"file":  filePath,
		"line":  line,
	}

	h.logger.Error("Frontend reguest error", fields)

	h.SetVar("panic", fields)
}
