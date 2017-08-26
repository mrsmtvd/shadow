package dashboard

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
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
