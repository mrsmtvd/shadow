package dashboard

import (
	"fmt"
	"net/http"
)

type PanicHandler struct {
	Handler
}

func (h *PanicHandler) ServeHTTP(w *Response, r *Request) {
	error := r.Panic()
	fields := map[string]interface{}{
		"error": fmt.Sprintf("%s", error.error),
		"stack": error.stack,
		"file":  error.file,
		"line":  error.line,
	}

	w.WriteHeader(http.StatusInternalServerError)
	h.RenderLayout(r.Context(), ComponentName, "500", "simple", map[string]interface{}{
		"panic": fields,
	})

	r.Logger().Error("Frontend reguest error", fields)
}
