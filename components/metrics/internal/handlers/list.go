package handlers

import (
	"sort"
	"time"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/metrics"
)

type ListHandler struct {
	dashboard.Handler
}

func (h *ListHandler) ServeHTTP(_ *dashboard.Response, r *dashboard.Request) {
	var updated time.Time

	measures, err := r.Component().(metrics.Component).Registry().Gather()
	if err == nil {
		sort.Sort(measures)

		for _, measure := range measures {
			if measure.CreatedAt.After(updated) {
				updated = measure.CreatedAt
			}
		}
	}

	h.Render(r.Context(), "list", map[string]interface{}{
		"measures": measures,
		"error":    err,
		"updated":  updated,
	})
}
