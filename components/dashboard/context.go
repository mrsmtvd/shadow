package dashboard

import (
	"context"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
)

var (
	ConfigContextKey   = &ContextKey{"config"}
	LoggerContextKey   = &ContextKey{"logger"}
	RenderContextKey   = &ContextKey{"render"}
	RequestContextKey  = &ContextKey{"request"}
	ResponseContextKey = &ContextKey{"response"}
	RouterContextKey   = &ContextKey{"router"}
	PanicContextKey    = &ContextKey{"panic"}
	SessionContextKey  = &ContextKey{"session"}
)

type ContextKey struct {
	name string
}

func (k *ContextKey) String() string {
	return "dashboard context value " + k.name
}

type PanicError struct {
	error interface{}
	stack string
	file  string
	line  int
}

func ConfigFromContext(c context.Context) *config.Component {
	return c.Value(ConfigContextKey).(*config.Component)
}

func LoggerFromContext(c context.Context) logger.Logger {
	return c.Value(LoggerContextKey).(logger.Logger)
}

func RenderFromContext(c context.Context) *Renderer {
	return c.Value(RenderContextKey).(*Renderer)
}

func RequestFromContext(c context.Context) *Request {
	return c.Value(RequestContextKey).(*Request)
}

func ResponseFromContext(c context.Context) *Response {
	return c.Value(ResponseContextKey).(*Response)
}

func RouterFromContext(c context.Context) *Router {
	return c.Value(RouterContextKey).(*Router)
}

func PanicFromContext(c context.Context) *PanicError {
	return c.Value(PanicContextKey).(*PanicError)
}

func SessionFromContext(c context.Context) *Session {
	return c.Value(SessionContextKey).(*Session)
}
