package frontend

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/rs/xlog"
)

type PanicHandler struct {
	AbstractFrontendHandler

	error interface{}
}

func (h *PanicHandler) SetError(err interface{}) {
	h.error = err
}

func (h *PanicHandler) Handle() {
	h.SetTemplate("500.tpl.html")
	h.SetPageTitle("Internal Server Error")

	h.Output.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.Output.WriteHeader(http.StatusInternalServerError)

	stack := make([]byte, 4096)
	stack = stack[:runtime.Stack(stack, false)]

	_, filePath, line, _ := runtime.Caller(0)

	fields := xlog.F{
		"error": fmt.Sprintf("%s", h.error),
		"stack": string(stack),
		"file":  filePath,
		"line":  line,
	}

	h.Service.(*FrontendService).logger.Error("Frontend reguest error", fields)

	h.SetVar("panic", fields)
}
