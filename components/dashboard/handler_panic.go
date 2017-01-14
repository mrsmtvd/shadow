package dashboard

import (
	"fmt"
	"net/http"
	"runtime"
)

type PanicHandler struct {
	Handler
}

func (h *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stack := make([]byte, 4096)
	stack = stack[:runtime.Stack(stack, false)]

	_, filePath, line, _ := runtime.Caller(0)

	fields := map[string]interface{}{
		"error": fmt.Sprintf("%s", PanicFromContext(r.Context())),
		"stack": string(stack),
		"file":  filePath,
		"line":  line,
	}

	w.WriteHeader(http.StatusInternalServerError)
	h.Render(r.Context(), "dashboard", "500", map[string]interface{}{
		"panic": fields,
	})

	LoggerFromContext(r.Context()).Error("Frontend reguest error", fields)
}
