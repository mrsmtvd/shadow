package dashboard

import (
	"context"
	"net/http"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
)

var (
	ConfigContextKey   = &ContextKey{"config"}
	LoggerContextKey   = &ContextKey{"logger"}
	RenderContextKey   = &ContextKey{"render"}
	RequestContextKey  = &ContextKey{"request"}
	ResponseContextKey = &ContextKey{"response"}
	PanicContextKey    = &ContextKey{"panic"}
)

type ContextKey struct {
	name string
}

func (k *ContextKey) String() string {
	return "dashboard context value " + k.name
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

func RequestFromContext(c context.Context) *http.Request {
	return c.Value(RequestContextKey).(*http.Request)
}

func ResponseFromContext(c context.Context) http.ResponseWriter {
	return c.Value(ResponseContextKey).(http.ResponseWriter)
}

func PanicFromContext(c context.Context) interface{} {
	return c.Value(PanicContextKey)
}
