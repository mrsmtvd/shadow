package handlers

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
)

type ChangeHandler struct {
	dashboard.Handler
}

func (h *ChangeHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	query := r.URL().Query().Get("locale")
	if query == "" {
		h.NotFound(w, r)
		return
	}

	component := r.Component().(i18n.Component)

	locale, ok := component.Manager().Locale(query)
	if !ok {
		h.NotFound(w, r)
		return
	}

	component.SaveToSession(r.Session(), locale)
	component.SaveToCookie(w, locale)

	redirect := r.URL().Query().Get("return")
	if redirect == "" {
		redirect = r.Original().Referer()
	}

	if redirect == "" || redirect == r.URL().Path {
		redirect = r.Config().String(dashboard.ConfigStartURL)

		if redirect == r.URL().Path {
			redirect = "/"
		}
	}

	h.Redirect(redirect, http.StatusTemporaryRedirect, w, r)
}
