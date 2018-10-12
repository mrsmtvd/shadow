package dashboard

import (
	"net/http"
)

type HasServerMiddleware interface {
	DashboardMiddleware() []func(http.Handler) http.Handler
}
