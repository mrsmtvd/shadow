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
	h.SetPageTitle("Configuration")
	h.SetPageHeader("Configuration")
	h.SetVar("Flags", flags)
}
