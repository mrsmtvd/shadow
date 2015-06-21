package frontend

import (
	"fmt"
	"net/http"
	"runtime"
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

	data := map[string]interface{}{
		"error": fmt.Sprintf("%s", h.error),
		"stack": string(stack),
		"file":  filePath,
		"line":  line,
	}

	h.View.Context["panic"] = data
}
