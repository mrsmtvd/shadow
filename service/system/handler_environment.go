package system

import (
	"runtime"

	"os"
	"strings"

	"github.com/kihamo/shadow/service/frontend"
)

type EnvironmentHandler struct {
	frontend.AbstractFrontendHandler
}

func (h *EnvironmentHandler) Handle() {
	h.SetTemplate("environment.tpl.html")
	h.SetPageTitle("Environment")
	h.SetPageHeader("Environment")

	vars := make(map[string]string, len(os.Environ()))
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		vars[pair[0]] = pair[1]
	}

	h.SetVar("GoVersion", runtime.Version())
	h.SetVar("GoOS", runtime.GOOS)
	h.SetVar("GoArch", runtime.GOARCH)
	h.SetVar("GoRoot", runtime.GOROOT())
	h.SetVar("EnvVars", vars)
}
