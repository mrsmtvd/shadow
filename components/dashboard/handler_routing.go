package dashboard

type RoutingHandler struct {
	Handler
}

func (h *RoutingHandler) ServeHTTP(_ *Response, r *Request) {
	router := RouterFromContext(r.Context())

	h.Render(r.Context(), ComponentName, "routing", map[string]interface{}{
		"routes": router.GetRoutes(),
	})
}
