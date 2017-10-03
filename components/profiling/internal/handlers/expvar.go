package handlers

import (
	"expvar"
	"fmt"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

type ExpvarHandler struct {
	dashboard.Handler
}

func (h *ExpvarHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if !r.Config().GetBool(config.ConfigDebug) {
		h.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n")
	first := true
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}
