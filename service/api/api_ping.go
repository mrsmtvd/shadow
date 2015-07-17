package api

import (
	"gopkg.in/jcelliott/turnpike.v2"
)

type PingProcedure struct {
	AbstractApiProcedure
}

func (p *PingProcedure) GetName() string {
	return "api.ping"
}

func (p *PingProcedure) Run([]interface{}, map[string]interface{}) *turnpike.CallResult {
	return &turnpike.CallResult{Args: []interface{}{"pong"}}
}
