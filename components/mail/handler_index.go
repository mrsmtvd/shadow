package mail

import (
	"github.com/kihamo/shadow/components/dashboard"
	"gopkg.in/gomail.v2"
)

type IndexHandler struct {
	dashboard.TemplateHandler

	component *Component
}

func (h *IndexHandler) Handle() {
	h.SetView("mail", "index")

	if h.IsPost() {
		message := gomail.NewMessage()
		message.SetHeader("Subject", h.Request().FormValue("subject"))
		message.SetHeader("To", h.Request().FormValue("to"))

		if h.Request().FormValue("type") == "html" {
			message.SetBody("text/html", h.Request().FormValue("message"))
		} else {
			message.SetBody("text/plain", h.Request().FormValue("message"))
		}

		if err := h.component.SendAndReturn(message); err != nil {
			h.SetVar("error", err.Error())
		} else {
			h.SetVar("message", "Message send success")
		}
	}
}
