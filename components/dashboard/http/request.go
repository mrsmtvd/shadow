package http

import (
	"encoding/json"
	"errors"
	originalHttp "net/http"

	"github.com/kihamo/gotypes"
)

type Request struct {
	original *originalHttp.Request
}

func NewRequest(r *originalHttp.Request) *Request {
	return &Request{
		original: r,
	}
}

func (r *Request) Original() *originalHttp.Request {
	return r.original
}

func (r *Request) IsGet() bool {
	return r.original.Method == originalHttp.MethodGet
}

func (r *Request) IsPost() bool {
	return r.original.Method == originalHttp.MethodPost
}

func (r *Request) IsAjax() bool {
	return r.original.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (r *Request) DecodeJSON(j interface{}) error {
	decoder := json.NewDecoder(r.original.Body)

	var in interface{}
	err := decoder.Decode(&in)

	if err != nil {
		return err
	}

	converter := gotypes.NewConverter(in, &j)

	if !converter.Valid() {
		return errors.New("Convert failed")
	}

	return nil
}
