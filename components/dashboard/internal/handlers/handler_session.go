package handlers

import (
	"net/http"

	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/dashboard"
)

type SessionHandler struct {
	dashboard.Handler
}

func (h *SessionHandler) ServeHTTP(w http.ResponseWriter, r *dashboard.Request) {
	if !r.Config().Bool(config.ConfigDebug) {
		h.NotFound(w, r)
		return
	}

	s := r.Session()

	keys := s.Keys()
	data := make(map[string]interface{}, len(keys))

	for _, key := range keys {
		data[key] = s.GetString(key)
	}

	h.Render(r.Context(), "session", map[string]interface{}{
		"data": data,
	})
}
