package grpc

import (
	"net/http"

	"github.com/kihamo/shadow/components/dashboard"
)

type IndexHandler struct {
	dashboard.Handler

	component *Component
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api := map[string][]string{}

	for service, info := range h.component.server.GetServiceInfo() {
		if _, ok := api[service]; !ok {
			api[service] = make([]string, 0)
		}

		for _, method := range info.Methods {
			api[service] = append(api[service], method.Name)
		}
	}

	h.Render(r.Context(), ComponentName, "index", map[string]interface{}{
		"api": api,
	})
}
