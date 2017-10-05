package handlers

import (
	"github.com/kihamo/shadow/components/dashboard"
)

type RoutingHandler struct {
	dashboard.Handler
}

func (h *RoutingHandler) ServeHTTP(_ *dashboard.Response, r *dashboard.Request) {
	router := dashboard.RouterFromContext(r.Context())

	h.Render(r.Context(), dashboard.ComponentName, "routing", map[string]interface{}{
		"routes": router.GetRoutes(),
	})
}
