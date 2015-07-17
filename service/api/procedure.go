package api

import (
	"github.com/kihamo/shadow"
	"gopkg.in/jcelliott/turnpike.v2"
)

type ApiProcedure interface {
	Init(shadow.Service, *shadow.Application)
	GetName() string
	Run([]interface{}, map[string]interface{}) *turnpike.CallResult
}

type AbstractApiProcedure struct {
	ApiProcedure
	Application *shadow.Application
	Service     shadow.Service
	ApiService  *ApiService
}

func (p *AbstractApiProcedure) Init(s shadow.Service, a *shadow.Application) {
	p.Application = a
	p.Service = s

	apiService, err := a.GetService("api")
	if err == nil {
		if castService, ok := apiService.(*ApiService); ok {
			p.ApiService = castService
			return
		}
	}

	panic("Api service not found")
}
