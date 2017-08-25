package http

import (
	"encoding/json"
	originalHttp "net/http"
)

type Response struct {
	originalHttp.ResponseWriter
	status int
}

func NewResponse(w originalHttp.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: w,
	}
}

func (w *Response) Header() originalHttp.Header {
	return w.ResponseWriter.Header()
}

func (w *Response) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = originalHttp.StatusOK
	}

	return w.ResponseWriter.Write(data)
}

func (w *Response) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *Response) GetStatusCode() int {
	return w.status
}

func (w *Response) SendJSON(r interface{}) []byte {
	response, err := json.Marshal(r)
	if err != nil {
		panic(err.Error())
	}

	w.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.ResponseWriter.Write(response)

	return response
}
