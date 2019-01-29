package handlers

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
)

type ChangeHandler struct {
	dashboard.Handler

	component i18n.Component
}

func NewChangeHandler(component i18n.Component) *ChangeHandler {
	return &ChangeHandler{
		component: component,
	}
}

func (h *ChangeHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	query := r.URL().Query().Get("locale")
	if query == "" {
		h.NotFound(w, r)
		return
	}

	locale, ok := h.component.Manager().Locale(query)
	if !ok {
		h.NotFound(w, r)
		return
	}

	if err := h.component.SaveToSession(r.Session(), locale); err != nil {
		panic(err.Error())
	}

	if err := h.component.SaveToCookie(w, locale); err != nil {
		panic(err.Error())
	}

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
