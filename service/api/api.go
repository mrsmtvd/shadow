package api

import (
	"gopkg.in/jcelliott/turnpike.v2"
)

func (s *ApiService) GetApiMethods() map[string]turnpike.MethodHandler {
	return map[string]turnpike.MethodHandler{
		"ping": s.Ping,
	}
}

func (s *ApiService) Ping([]interface{}, map[string]interface{}) *turnpike.CallResult {
	return &turnpike.CallResult{Args: []interface{}{"pong"}}
}
