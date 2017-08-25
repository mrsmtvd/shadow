package dashboard

import (
	"context"
	"net/http"
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

func (h *Handler) Redirect(l string, c int, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, l, c)
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	RouterFromContext(r.Context()).NotFound.ServeHTTP(w, r)
}

func (h *Handler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	RouterFromContext(r.Context()).MethodNotAllowed.ServeHTTP(w, r)
}

func (h *Handler) Render(ctx context.Context, c, v string, d map[string]interface{}) {
	err := RenderFromContext(ctx).Render(ctx, c, v, d)

	if err != nil {
		panic(err.Error())
	}
}

func (h *Handler) RenderLayout(ctx context.Context, c, v, l string, d map[string]interface{}) {
	err := RenderFromContext(ctx).RenderLayout(ctx, c, v, l, d)

	if err != nil {
		panic(err.Error())
	}
}
