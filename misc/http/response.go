package http

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
	length int64
	status int
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: w,
	}
}

func (w *Response) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *Response) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}

	n, err := w.ResponseWriter.Write(data)
	w.length += int64(n)

	return n, err
}

func (w *Response) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *Response) StatusCode() int {
	return w.status
}

func (w *Response) Length() int64 {
	return w.length
}

func (w *Response) SendJSON(r interface{}) error {
	response, err := json.Marshal(r)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(response)

	return err
}
