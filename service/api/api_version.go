package api

import (
	"gopkg.in/jcelliott/turnpike.v2"
)

type VersionProcedure struct {
	AbstractApiProcedure
}

func (p *VersionProcedure) GetName() string {
	return "api.version"
}

func (p *VersionProcedure) Run([]interface{}, map[string]interface{}) *turnpike.CallResult {
	return &turnpike.CallResult{Kwargs: map[string]interface{}{
		"version": p.Application.Version,
		"build":   p.Application.Build,
	}}
}
