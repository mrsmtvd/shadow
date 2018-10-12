package dashboard

import (
	o "net/http"

	"github.com/kihamo/shadow/misc/http"
)

type Response = http.Response

func NewResponse(w o.ResponseWriter) *Response {
	return http.NewResponse(w)
}
