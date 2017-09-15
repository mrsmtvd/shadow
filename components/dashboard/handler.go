package dashboard

import (
	"context"
	"net/http"
)

type Handler struct {
	http.Handler
}

func FromRouteHandler(h RouterHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		h.ServeHTTP(NewResponse(w), NewRequest(rq))
	})
}

func (h *Handler) Redirect(l string, c int, w *Response, r *Request) {
	http.Redirect(w, r.Original(), l, c)
}

func (h *Handler) NotFound(w *Response, r *Request) {
	router := RouterFromContext(r.Context())
	if router == nil {
		panic("Router isn't set in context")
	}

	router.NotFound.ServeHTTP(w, r.Original())
}

func (h *Handler) MethodNotAllowed(w *Response, r *Request) {
	router := RouterFromContext(r.Context())
	if router == nil {
		panic("Router isn't set in context")
	}

	router.MethodNotAllowed.ServeHTTP(w, r.Original())
}

func (h *Handler) Render(ctx context.Context, c, v string, d map[string]interface{}) {
	render := RenderFromContext(ctx)
	if render == nil {
		panic("Render isn't set in context")
	}

	if err := render.Render(ctx, c, v, d); err != nil {
		panic(err.Error())
	}
}

func (h *Handler) RenderLayout(ctx context.Context, c, v, l string, d map[string]interface{}) {
	render := RenderFromContext(ctx)
	if render == nil {
		panic("Render isn't set in context")
	}

	if err := render.RenderLayout(ctx, c, v, l, d); err != nil {
		panic(err.Error())
	}
}
