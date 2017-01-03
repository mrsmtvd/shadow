package dashboard

import (
	"os"
	"runtime"
	"strings"
)

type EnvironmentHandler struct {
	TemplateHandler
}

func (h *EnvironmentHandler) Handle() {
	h.SetView("dashboard", "environment")

	vars := make(map[string]string, len(os.Environ()))
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		vars[pair[0]] = pair[1]
	}

	h.SetVar("goVersion", runtime.Version())
	h.SetVar("goOS", runtime.GOOS)
	h.SetVar("goArch", runtime.GOARCH)
	h.SetVar("goRoot", runtime.GOROOT())
	h.SetVar("goNumCPU", runtime.NumCPU())
	h.SetVar("goNumGoroutine", runtime.NumGoroutine())
	h.SetVar("goNumCgoCall", runtime.NumCgoCall())
	h.SetVar("envVars", vars)
}
