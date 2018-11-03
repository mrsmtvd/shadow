package dashboard

import (
	"context"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
)

var (
	ApplicationContextKey = &ContextKey{"application"}
	ComponentContextKey   = &ContextKey{"component"}
	ConfigContextKey      = &ContextKey{"config"}
	PanicContextKey       = &ContextKey{"panic"}
	RenderContextKey      = &ContextKey{"render"}
	RequestContextKey     = &ContextKey{"request"}
	ResponseContextKey    = &ContextKey{"response"}
	RouteContextKey       = &ContextKey{"route"}
	RouterContextKey      = &ContextKey{"router"}
	SessionContextKey     = &ContextKey{"session"}
)

type ContextKey struct {
	Name string
}

func (k *ContextKey) String() string {
	return "dashboard context value " + k.Name
}

func ApplicationFromContext(c context.Context) shadow.Application {
	v := c.Value(ApplicationContextKey)
	if v != nil {
		return v.(shadow.Application)
	}

	return nil
}

func ComponentFromContext(c context.Context) shadow.Component {
	v := c.Value(ComponentContextKey)
	if v != nil {
		return v.(shadow.Component)
	}

	return nil
}

func ConfigFromContext(c context.Context) config.Component {
	v := c.Value(ConfigContextKey)
	if v != nil {
		return v.(config.Component)
	}

	return nil
}

func PanicFromContext(c context.Context) *PanicError {
	v := c.Value(PanicContextKey)
	if v != nil {
		return v.(*PanicError)
	}

	return nil
}

func RenderFromContext(c context.Context) Renderer {
	v := c.Value(RenderContextKey)
	if v != nil {
		return v.(Renderer)
	}

	return nil
}

func RequestFromContext(c context.Context) *Request {
	v := c.Value(RequestContextKey)
	if v != nil {
		return v.(*Request)
	}

	return nil
}

func ResponseFromContext(c context.Context) *Response {
	v := c.Value(ResponseContextKey)
	if v != nil {
		return v.(*Response)
	}

	return nil
}

func RouteFromContext(c context.Context) Route {
	v := c.Value(RouteContextKey)
	if v != nil {
		return v.(Route)
	}

	return nil
}

func RouterFromContext(c context.Context) Router {
	v := c.Value(RouterContextKey)
	if v != nil {
		return v.(Router)
	}

	return nil
}

func SessionFromContext(c context.Context) Session {
	v := c.Value(SessionContextKey)
	if v != nil {
		return v.(Session)
	}

	return nil
}
