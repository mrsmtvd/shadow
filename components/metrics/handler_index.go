package metrics

import (
	"net/http"
	"sort"
	"time"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/snitch"
)

type IndextHandler struct {
	dashboard.Handler
}

func (h *IndextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var updated time.Time

	metrics, err := snitch.DefaultRegisterer.Gather()
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
