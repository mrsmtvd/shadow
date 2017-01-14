package workers

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
)

type IndexHandler struct {
	dashboard.Handler
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Render(r.Context(), "workers", "index", nil)
}
