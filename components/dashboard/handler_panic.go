package dashboard

import (
	"fmt"
	"net/http"
)

type PanicHandler struct {
	Handler
}

func (h *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	error := PanicFromContext(r.Context())
	fields := map[string]interface{}{
		"error": fmt.Sprintf("%s", error.error),
		"stack": error.stack,
		"file":  error.file,
		"line":  error.line,
	}

	w.WriteHeader(http.StatusInternalServerError)
	h.Render(r.Context(), "dashboard", "500", map[string]interface{}{
		"panic": fields,
	})

	LoggerFromContext(r.Context()).Error("Frontend reguest error", fields)
}
