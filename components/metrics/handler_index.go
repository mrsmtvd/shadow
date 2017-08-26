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

	metrics, err := h.component.Registry().Gather()
	if err == nil {
		sort.Sort(metrics)

		for _, metric := range metrics {
			if metric.CreatedAt.After(updated) {
				updated = metric.CreatedAt
			}
		}
	}

	h.Render(r.Context(), ComponentName, "list", map[string]interface{}{
		"metrics": metrics,
		"error":   err,
		"updated": updated,
	})
}
