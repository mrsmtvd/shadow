package handlers

import (
	"net/http"

	"github.com/heptiolabs/healthcheck"
	"github.com/kihamo/snitch"
	"github.com/mrsmtvd/shadow"
	"github.com/mrsmtvd/shadow/components/dashboard"
)

type HealthCheckHandler struct {
	dashboard.Handler

	healthCheck             healthcheck.Handler
	metricHealthCheckStatus snitch.Gauge
}

func NewHealthCheckHandler(components []shadow.Component, m snitch.Gauge) *HealthCheckHandler {
	h := &HealthCheckHandler{
		healthCheck:             healthcheck.NewHandler(),
		metricHealthCheckStatus: m,
	}

	for _, component := range components {
		if componentLivenessCheck, ok := component.(dashboard.HasLivenessCheck); ok {
			for name, check := range componentLivenessCheck.LivenessCheck() {
				checkName := h.getCheckName(component, name)
				h.healthCheck.AddLivenessCheck(checkName, h.wrap(checkName, check))
			}
		}

		if componentReadinessCheck, ok := component.(dashboard.HasReadinessCheck); ok {
			for name, check := range componentReadinessCheck.ReadinessCheck() {
				checkName := h.getCheckName(component, name)
				h.healthCheck.AddReadinessCheck(checkName, h.wrap(checkName, check))
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

func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *dashboard.Request) {
	switch r.URL().Query().Get(":healthcheck") {
	case "live":
		h.healthCheck.LiveEndpoint(w, r.Original())
		return

	case "ready":
		h.healthCheck.ReadyEndpoint(w, r.Original())
		return
	}

	h.NotFound(w, r)
}

func (h *HealthCheckHandler) wrap(name string, check dashboard.HealthCheck) dashboard.HealthCheck {
	return func() (err error) {
		err = check()

		if err == nil {
			h.metricHealthCheckStatus.With("check", name).Set(0)
		} else {
			h.metricHealthCheckStatus.With("check", name).Set(1)
		}

		return err
	}
}
