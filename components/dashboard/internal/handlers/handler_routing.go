package handlers

import (
	"github.com/kihamo/shadow/components/dashboard"
)

type RoutingHandler struct {
	dashboard.Handler

	router dashboard.Router
}

func NewRoutingHandler(router dashboard.Router) *RoutingHandler {
	return &RoutingHandler{
		router: router,
	}
}

func (h *RoutingHandler) ServeHTTP(_ *dashboard.Response, r *dashboard.Request) {
	h.Render(r.Context(), "routing", map[string]interface{}{
		"routes": h.router.Routes(),
	})
}
