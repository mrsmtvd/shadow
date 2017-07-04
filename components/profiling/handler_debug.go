package profiling

import (
	"net/http"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) debugHandler(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c.config.GetBool(config.ConfigDebug) {
			h.ServeHTTP(w, r)
		} else {
			dashboard.RouterFromContext(r.Context()).NotFound.ServeHTTP(w, r)
		}
	})
}
