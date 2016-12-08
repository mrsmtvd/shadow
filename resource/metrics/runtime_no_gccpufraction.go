// +build !go1.5

package metrics

import (
	"runtime"
)

func gcCPUFraction(_ *runtime.MemStats) float64 {
	return 0
}
