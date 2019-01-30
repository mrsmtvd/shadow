package dashboard

import (
	"context"
)

type contextKey string

var (
	templateNamespaceContextKey = contextKey("template-namespace")
	panicContextKey             = contextKey("panic")
	renderContextKey            = contextKey("render")
	requestContextKey           = contextKey("request")
	responseContextKey          = contextKey("response")
	routeContextKey             = contextKey("route")
	routerContextKey            = contextKey("router")
	sessionContextKey           = contextKey("session")
)

func ContextWithTemplateNamespace(ctx context.Context, ns string) context.Context {
	return context.WithValue(ctx, templateNamespaceContextKey, ns)
}

func TemplateNamespaceFromContext(c context.Context) string {
	v := c.Value(templateNamespaceContextKey)
	if v != nil {
		return v.(string)
	}

	return ComponentName
}

func ContextWithPanic(ctx context.Context, err *PanicError) context.Context {
	return context.WithValue(ctx, panicContextKey, err)
}

func PanicFromContext(c context.Context) *PanicError {
	v := c.Value(panicContextKey)
	if v != nil {
		return v.(*PanicError)
	}

	return nil
}

func ContextWithRender(ctx context.Context, render Renderer) context.Context {
	return context.WithValue(ctx, renderContextKey, render)
}

func RenderFromContext(c context.Context) Renderer {
	v := c.Value(renderContextKey)
	if v != nil {
		return v.(Renderer)
	}

	return nil
}

func ContextWithRequest(ctx context.Context, request *Request) context.Context {
	return context.WithValue(ctx, requestContextKey, request)
}

func RequestFromContext(c context.Context) *Request {
	v := c.Value(requestContextKey)
	if v != nil {
		return v.(*Request)
	}

	return nil
}

func ContextWithResponse(ctx context.Context, response *Response) context.Context {
	return context.WithValue(ctx, responseContextKey, response)
}

func ResponseFromContext(c context.Context) *Response {
	v := c.Value(responseContextKey)
	if v != nil {
		return v.(*Response)
	}

	return nil
}

func ContextWithRoute(ctx context.Context, route Route) context.Context {
	return context.WithValue(ctx, routeContextKey, route)
}

func RouteFromContext(c context.Context) Route {
	v := c.Value(routeContextKey)
	if v != nil {
		return v.(Route)
	}

	return nil
}

func ContextWithRouter(ctx context.Context, router Router) context.Context {
	return context.WithValue(ctx, routerContextKey, router)
}

func RouterFromContext(c context.Context) Router {
	v := c.Value(routerContextKey)
	if v != nil {
		return v.(Router)
	}

	return nil
}

func ContextWithSession(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, session)
}

func SessionFromContext(c context.Context) Session {
	v := c.Value(sessionContextKey)
	if v != nil {
		return v.(Session)
	}

	return nil
}
