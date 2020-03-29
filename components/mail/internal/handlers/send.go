package handlers

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/mail"
	"gopkg.in/gomail.v2"
)

type SendHandler struct {
	dashboard.Handler

	component mail.Component
}

func NewSendHandler(component mail.Component) *SendHandler {
	return &SendHandler{
		component: component,
	}
}

func (h *SendHandler) ServeHTTP(w http.ResponseWriter, r *dashboard.Request) {
	locale := i18n.Locale(r.Context())
	vars := map[string]interface{}{}

	if r.IsPost() {
		message := gomail.NewMessage()
		message.SetHeader("Subject", r.Original().FormValue("subject"))
		message.SetHeader("To", r.Original().FormValue("to"))

		if r.Original().FormValue("type") == "html" {
			message.SetBody("text/html", r.Original().FormValue("message"))
		} else {
			message.SetBody("text/plain", r.Original().FormValue("message"))
		}

		if err := h.component.SendAndReturn(message); err != nil {
			r.Session().FlashBag().Error(err.Error())
		} else {
			r.Session().FlashBag().Success(locale.Translate(h.component.Name(), "Message send success", ""))
			h.Redirect(r.URL().String(), http.StatusFound, w, r)
			return
		}
	}

	h.Render(r.Context(), "send", vars)
}
