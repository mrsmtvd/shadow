package api

import (
	"fmt"
	"os"

	"github.com/kihamo/shadow/resource"
)

func (s *ApiService) GetConfigVariables() []resource.ConfigVariable {
	pathCrt := ""
	pathKey := ""

	dir, err := os.Getwd()
	if err == nil {
		pathCrt = fmt.Sprint(dir, "/server.crt")
		pathKey = fmt.Sprint(dir, "/server.key")
	}

	return []resource.ConfigVariable{
		resource.ConfigVariable{
			Key:   "api-host",
			Value: "0.0.0.0",
			Usage: "API socket host",
		},
		resource.ConfigVariable{
			Key:   "api-port",
			Value: 8001,
			Usage: "API socket port",
		},
		resource.ConfigVariable{
			Key:   "api-secure",
			Value: false,
			Usage: "API enable SSL",
		},
		resource.ConfigVariable{
			Key:   "api-secure-crt",
			Value: pathCrt,
			Usage: "API path to SSL crt file",
		},
		resource.ConfigVariable{
			Key:   "api-secure-key",
			Value: pathKey,
			Usage: "API path to SSL key file",
		},
	}
}
