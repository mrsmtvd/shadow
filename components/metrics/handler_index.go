package metrics

import (
	"sort"
	"time"

	"github.com/kihamo/shadow/components/dashboard"
)

type IndextHandler struct {
	dashboard.Handler

	component *Component
}

func (h *IndextHandler) ServeHTTP(_ *dashboard.Response, r *dashboard.Request) {
	var updated time.Time

	measures, err := h.component.Registry().Gather()
	if err == nil {
		sort.Sort(measures)

		for _, measure := range measures {
			if measure.CreatedAt.After(updated) {
				updated = measure.CreatedAt
			}
		}
	}

	h.Render(r.Context(), ComponentName, "list", map[string]interface{}{
		"measures": measures,
		"error":    err,
		"updated":  updated,
	})
}
