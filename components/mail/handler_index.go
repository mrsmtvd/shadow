package mail

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	"gopkg.in/gomail.v2"
)

type IndexHandler struct {
	dashboard.Handler

	component *Component
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := map[string]interface{}{}
	request := dashboard.RequestFromContext(r.Context())

	if request.IsPost() {
		message := gomail.NewMessage()
		message.SetHeader("Subject", r.FormValue("subject"))
		message.SetHeader("To", r.FormValue("to"))

		if r.FormValue("type") == "html" {
			message.SetBody("text/html", r.FormValue("message"))
		} else {
			message.SetBody("text/plain", r.FormValue("message"))
		}

		if err := h.component.SendAndReturn(message); err != nil {
			vars["error"] = err.Error()
		} else {
			vars["message"] = "Message send success"
		}
	}

	h.Render(r.Context(), ComponentName, "index", vars)
}
