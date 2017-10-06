package handlers

import (
	"net/http"

	"github.com/kihamo/shadow/components/alerts"
	"github.com/kihamo/shadow/components/dashboard"
)

type SendHandler struct {
	dashboard.Handler

	Component alerts.Component
}

func (h *SendHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if r.IsPost() {
		h.Component.Send(
			r.Original().FormValue("title"),
			r.Original().FormValue("message"),
			r.Original().FormValue("icon"))

		h.Redirect(dashboard.RouteFromContext(r.Context()).Path(), http.StatusFound, w, r)
		return
	}

	h.Render(r.Context(), h.Component.GetName(), "send", nil)
}
