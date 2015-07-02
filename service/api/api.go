package api

import (
	base "github.com/apache/thrift/lib/go/thrift"
	gen "github.com/kihamo/shadow/service/api/gen-go/api"
)

func (s *ApiService) GetProcessor() base.TProcessor {
	return gen.NewApiProcessor(&ApiHandler{})
}

type ApiHandler struct {
}

func (h *ApiHandler) Ping() (bool, error) {
	return true, nil
}
