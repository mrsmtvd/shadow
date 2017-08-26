package mail

import (
	"github.com/kihamo/shadow/components/dashboard"
	"gopkg.in/gomail.v2"
)

type IndexHandler struct {
	dashboard.Handler

	component *Component
}

func (h *IndexHandler) ServeHTTP(_ *dashboard.Response, r *dashboard.Request) {
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
			vars["error"] = err.Error()
		} else {
			vars["message"] = "Message send success"
		}
	}

	h.Render(r.Context(), ComponentName, "index", vars)
}
