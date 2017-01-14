package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kihamo/gotypes"
)

type HandlerAuth interface {
	IsAuth() bool
}

type Handler struct {
	http.Handler
}

func (h *Handler) IsAuth() bool {
	return true
}

func (h *Handler) IsGet(r *http.Request) bool {
	return r.Method == http.MethodGet
}

func (h *Handler) IsPost(r *http.Request) bool {
	return r.Method == http.MethodPost
}

func (h *Handler) IsAjax(r *http.Request) bool {
	return r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (h *Handler) Redirect(l string, c int, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, l, c)
}

func (h *Handler) Render(ctx context.Context, c, v string, d map[string]interface{}) {
	err := RenderFromContext(ctx).Render(ctx, c, v, d)

	if err != nil {
		panic(err.Error())
	}
}

func (h *Handler) SendJSON(r interface{}, w http.ResponseWriter) {
	response, err := json.Marshal(r)
	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(response)
}

func (h *Handler) DecodeJSON(j interface{}, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)

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
