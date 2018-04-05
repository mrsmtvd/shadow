package handlers

import (
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/mail"
	"gopkg.in/gomail.v2"
)

type SendHandler struct {
	dashboard.Handler
}

func (h *SendHandler) ServeHTTP(_ *dashboard.Response, r *dashboard.Request) {
	locale := i18n.NewOrNopFromRequest(r)
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

		if err := r.Component().(mail.Component).SendAndReturn(message); err != nil {
			vars["error"] = err.Error()
		} else {
			vars["message"] = locale.Translate(r.Component().Name(), "Message send success", "")
		}
	}

	h.Render(r.Context(), "send", vars)
}
