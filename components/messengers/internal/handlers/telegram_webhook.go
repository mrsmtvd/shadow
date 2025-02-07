package handlers

import (
	"net/http"

	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/messengers"
)

type TelegramWebHookHandler struct {
	dashboard.Handler
}

func (h *TelegramWebHookHandler) ServeHTTP(w http.ResponseWriter, r *dashboard.Request) {
	if !r.Config().Bool(messengers.ConfigTelegramWebHookEnabled) || !r.Config().Bool(messengers.ConfigTelegramUpdatesEnabled) {
		h.NotFound(w, r)
		return
	}
}
