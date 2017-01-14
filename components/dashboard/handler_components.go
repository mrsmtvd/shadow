package dashboard

import (
	"net/http"

	"github.com/kihamo/shadow"
)

type ComponentsHandler struct {
	Handler

	application shadow.Application
}

func (h *ComponentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Render(r.Context(), "dashboard", "components", map[string]interface{}{
		"components": h.application.GetComponents(),
	})
}
