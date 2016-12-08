package metrics

import (
	"time"
)

func CaptureMetrics(d time.Duration, f func()) {
	go func() {
		for range time.NewTicker(d).C {
			f()
		}
	}()
}
