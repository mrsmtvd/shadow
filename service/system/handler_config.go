package system

import (
	"flag"

	"github.com/kihamo/shadow/service/frontend"
)

type ConfigHandler struct {
	frontend.AbstractFrontendHandler
}

func (h *ConfigHandler) Handle() {
	flags := []*flag.Flag{}
	flag.VisitAll(func(flag *flag.Flag) {
		flags = append(flags, flag)
	})

	h.SetTemplate("config.tpl.html")
	h.View.Context["PageTitle"] = "Configuration"
	h.View.Context["Flags"] = flags
}
