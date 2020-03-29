package handlers

import (
	"os"
	"runtime"
	"strings"

	"github.com/kihamo/shadow/components/dashboard"
	rt "github.com/kihamo/shadow/components/dashboard/runtime"
)

type EnvironmentHandler struct {
	dashboard.Handler
}

func (h *EnvironmentHandler) ServeHTTP(_ *dashboard.Response, r *dashboard.Request) {
	vars := make(map[string]string, len(os.Environ()))

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		vars[pair[0]] = pair[1]
	}

	h.Render(r.Context(), "environment", map[string]interface{}{
		"goVersion":      runtime.Version(),
		"goOS":           runtime.GOOS,
		"goArch":         runtime.GOARCH,
		"goRoot":         runtime.GOROOT(),
		"goNumCPU":       runtime.NumCPU(),
		"goNumGoroutine": runtime.NumGoroutine(),
		"goNumCgoCall":   runtime.NumCgoCall(),
		"goRace":         rt.RaceEnabled,
		"envVars":        vars,
	})
}
