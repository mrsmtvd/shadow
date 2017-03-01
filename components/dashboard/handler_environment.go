package dashboard

import (
	"net/http"
	"os"
	"runtime"
	"strings"
)

type EnvironmentHandler struct {
	Handler
}

func (h *EnvironmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := make(map[string]string, len(os.Environ()))
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		vars[pair[0]] = pair[1]
	}

	h.Render(r.Context(), ComponentName, "environment", map[string]interface{}{
		"goVersion":      runtime.Version(),
		"goOS":           runtime.GOOS,
		"goArch":         runtime.GOARCH,
		"goRoot":         runtime.GOROOT(),
		"goNumCPU":       runtime.NumCPU(),
		"goNumGoroutine": runtime.NumGoroutine(),
		"goNumCgoCall":   runtime.NumCgoCall(),
		"envVars":        vars,
	})
}
