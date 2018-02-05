package handlers

import (
	"fmt"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/messengers"
)

type TelegramWebHookHandler struct {
	dashboard.Handler
}

func (h *TelegramWebHookHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if !r.Config().Bool(messengers.ConfigTelegramWebHookEnabled) || !r.Config().Bool(messengers.ConfigTelegramUpdatesEnabled) {
		h.NotFound(w, r)
		return
	}

	fmt.Println("Callme")
}
