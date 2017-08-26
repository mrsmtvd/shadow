package profiling

import (
	"expvar"
	"fmt"

	"github.com/kihamo/shadow/components/dashboard"
)

type ExpvarHandler struct {
}

func (h *ExpvarHandler) ServeHTTP(w *dashboard.Response, _ *dashboard.Request) {
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
