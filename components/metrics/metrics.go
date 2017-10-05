package metrics

import (
	"github.com/kihamo/snitch"
)

type HasMetrics interface {
	Metrics() snitch.Collector
}
