package api

import (
	"github.com/kihamo/shadow/resource"
)

func (s *ApiService) GetConfigVariables() []resource.ConfigVariable {
	return []resource.ConfigVariable{
		resource.ConfigVariable{
			Key:   "api-host",
			Value: "0.0.0.0",
			Usage: "API socket host",
		},
		resource.ConfigVariable{
			Key:   "api-port",
			Value: "8001",
			Usage: "API socket port",
		},
		resource.ConfigVariable{
			Key:   "api-protocol",
			Value: "binary",
			Usage: "API protocol: binary, compact, json, simplejson",
		},
		resource.ConfigVariable{
			Key:   "api-transport",
			Value: "framed",
			Usage: "API transport: buffered, framed, \"\"",
		},
	}
}
