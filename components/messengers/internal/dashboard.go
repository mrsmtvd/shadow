package internal

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/messengers/internal/handlers"
)

func (c *Component) DashboardRoutes() []dashboard.Route {
	return []dashboard.Route{
		dashboard.NewRoute(
			c.Name(),
			[]string{http.MethodGet},
			"/"+c.Name()+"/telegram/webhook",
			&handlers.TelegramWebHookHandler{},
			"",
			true),
	}
}
