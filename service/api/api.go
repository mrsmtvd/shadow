package api

import (
	"gopkg.in/jcelliott/turnpike.v2"
)

func (s *ApiService) GetApiMethods() map[string]turnpike.MethodHandler {
	return map[string]turnpike.MethodHandler{
		"ping": s.Ping,
        "version": s.Version,
	}
}

func (s *ApiService) Ping([]interface{}, map[string]interface{}) *turnpike.CallResult {
	return &turnpike.CallResult{Args: []interface{}{"pong"}}
}

func (s *ApiService) Version([]interface{}, map[string]interface{}) *turnpike.CallResult {
    return &turnpike.CallResult{Kwargs: map[string]interface{}{
        "version": s.application.Version,
        "build": s.application.Build,
    }}
}
