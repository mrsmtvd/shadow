package dashboard

import (
	"net/http"
)

type Router interface {
	Routes() []Route
	NotFoundServeHTTP(w http.ResponseWriter, r *http.Request)
	MethodNotAllowedServeHTTP(w http.ResponseWriter, r *http.Request)
	InternalErrorServeHTTP(w http.ResponseWriter, r *http.Request, e error)
}
