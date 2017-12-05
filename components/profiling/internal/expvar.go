package internal

import (
	"runtime"
)

func expvarRuntime() interface{} {
	return map[string]interface{}{
		"gomaxproc": runtime.GOMAXPROCS(-1),
		"numcpu":    runtime.NumCPU(),
	}
}
