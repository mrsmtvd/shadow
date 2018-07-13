package dashboard

import (
	"github.com/heptiolabs/healthcheck"
)

type HealthCheck = healthcheck.Check

type HasLivenessCheck interface {
	LivenessCheck() map[string]HealthCheck
}

type HasReadinessCheck interface {
	ReadinessCheck() map[string]HealthCheck
}
