package handlers

import (
	"github.com/heptiolabs/healthcheck"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/dashboard"
)

type HealthCheckHandler struct {
	dashboard.Handler

	healthCheck healthcheck.Handler
	hasCheck    bool
}

func NewHealthCheckHandler(a shadow.Application) *HealthCheckHandler {
	h := &HealthCheckHandler{
		healthCheck: healthcheck.NewHandler(),
	}

	components, err := a.GetComponents()
	if err != nil {
		return h
	}

	for _, component := range components {
		if componentLivenessCheck, ok := component.(dashboard.HasLivenessCheck); ok {
			for name, check := range componentLivenessCheck.LivenessCheck() {
				h.healthCheck.AddLivenessCheck(h.getCheckName(component, name), check)
				h.hasCheck = true
			}
		}

		if componentReadinessCheck, ok := component.(dashboard.HasReadinessCheck); ok {
			for name, check := range componentReadinessCheck.ReadinessCheck() {
				h.healthCheck.AddReadinessCheck(h.getCheckName(component, name), check)
				h.hasCheck = true
			}
		}
	}

	return h
}

func (h *HealthCheckHandler) getCheckName(cmp shadow.Component, name string) string {
	if name != "" {
		return cmp.Name() + "_" + name
	}

	return cmp.Name()
}

func (h *HealthCheckHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if h.hasCheck {
		switch r.URL().Query().Get(":healthcheck") {
		case "live":
			h.healthCheck.LiveEndpoint(w, r.Original())

		case "ready":
			h.healthCheck.ReadyEndpoint(w, r.Original())
		}

		return
	}

	h.NotFound(w, r)
}
