package dashboard

import (
	"net/http"
)

type Router interface {
	GetRoutes() []Route
	NotFoundServeHTTP(w http.ResponseWriter, r *http.Request)
	MethodNotAllowedServeHTTP(w http.ResponseWriter, r *http.Request)
}
