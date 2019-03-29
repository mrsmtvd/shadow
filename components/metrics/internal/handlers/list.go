package handlers

import (
	"sort"
	"time"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/metrics"
)

type ListHandler struct {
	dashboard.Handler

	component metrics.Component
}

func NewListHandler(component metrics.Component) *ListHandler {
	return &ListHandler{
		component: component,
	}
}

func (h *ListHandler) ServeHTTP(_ *dashboard.Response, r *dashboard.Request) {
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

	if err != nil {
		r.Session().FlashBag().Error(err.Error())
	}

	h.Render(r.Context(), "list", map[string]interface{}{
		"measures": measures,
		"updated":  updated,
	})
}
