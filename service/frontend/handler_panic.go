package frontend

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/Sirupsen/logrus"
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
	h.View.Context["PageTitle"] = "Internal Server Error"

	h.Output.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.Output.WriteHeader(http.StatusInternalServerError)

	stack := make([]byte, 4096)
	stack = stack[:runtime.Stack(stack, false)]

	_, filePath, line, _ := runtime.Caller(0)

	fields := logrus.Fields{
		"error": fmt.Sprintf("%s", h.error),
		"stack": string(stack),
		"file":  filePath,
		"line":  line,
	}

	h.Service.(*FrontendService).Logger.
		WithFields(fields).
		Error("Frontend reguest error")

	h.View.Context["panic"] = fields
}
