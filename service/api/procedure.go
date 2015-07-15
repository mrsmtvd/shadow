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

func (c *AbstractApiProcedure) Init(s shadow.Service, a *shadow.Application) {
	c.Application = a
	c.Service = s

	apiService, err := a.GetService("api")
	if err == nil {
		if castService, ok := apiService.(*ApiService); ok {
			c.ApiService = castService
			return
		}
	}

	panic("Api service not found")
}
