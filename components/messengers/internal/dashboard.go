package internal

import (
	"net/http"

	"github.com/mrsmtvd/shadow/components/dashboard"
	"github.com/mrsmtvd/shadow/components/messengers/internal/handlers"
)

func (c *Component) DashboardRoutes() []dashboard.Route {
	return []dashboard.Route{
		dashboard.NewRoute("/"+c.Name()+"/telegram/webhook", &handlers.TelegramWebHookHandler{}).
			WithMethods([]string{http.MethodGet}).
			WithAuth(true),
	}
}
