package internal

import (
	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/i18n"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(i18n.ConfigLocaleCookieName, config.ValueTypeString).
			WithUsage("Cookie name").
			WithGroup("Locale save").
			WithEditable(true).
			WithDefault("locale"),
		config.NewVariable(i18n.ConfigLocaleSessionKey, config.ValueTypeString).
			WithUsage("Key in session").
			WithGroup("Locale save").
			WithEditable(true).
			WithDefault("locale"),
	}
}
