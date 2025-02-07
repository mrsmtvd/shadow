package internal

import (
	"os"

	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/ota"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(ota.ConfigReleasesDirectory, config.ValueTypeString).
			WithUsage("Path to saved releases directory").
			WithEditable(true).
			WithDefault(os.TempDir()),
		config.NewVariable(ota.ConfigRepositoryServerEnabled, config.ValueTypeBool).
			WithUsage("Enable serve repository").
			WithGroup("Server").
			WithEditable(true).
			WithDefault(true),
		config.NewVariable(ota.ConfigRepositoryClientShadow, config.ValueTypeString).
			WithUsage("Shadow").
			WithGroup("Clients").
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add a url"}),
	}
}
