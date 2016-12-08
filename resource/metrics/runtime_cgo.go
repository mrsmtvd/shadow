// +build cgo
// +build !appengine

package metrics

import (
	"runtime"
)

func getNumCgoCall() int64 {
	return runtime.NumCgoCall()
}
