package frontend

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
)

type Router struct {
	httprouter.Router
	application *shadow.Application

	loggerMiddleware alice.Constructor
	authMiddleware   alice.Constructor
}

func NewRouter(application *shadow.Application, logger *logrus.Entry, config *resource.Config) *Router {
	r := &Router{}
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true
	r.HandleMethodNotAllowed = true

	r.application = application
	r.loggerMiddleware = LoggerMiddleware(logger)
	r.authMiddleware = BasicAuthMiddleware(config.GetString("frontend.auth-user"), config.GetString("frontend.auth-password"))

	return r
}

func (r *Router) GET(s shadow.Service, path string, h interface{}) {
	r.Handle(s, "GET", path, h)
}

func (r *Router) POST(s shadow.Service, path string, h interface{}) {
	r.Handle(s, "POST", path, h)
}

func (r *Router) Handle(s shadow.Service, m, p string, h interface{}) {
	var chain alice.Chain

	if h1, ok := h.(Handler); ok {
		chain = alice.New(
			r.loggerMiddleware,
		)

		if auth, ok := h1.(HandlerAuth); ok && auth.IsAuth() {
			chain = chain.Append(r.authMiddleware)
		}

		r.Router.Handler(m, p, chain.ThenFunc(func(out http.ResponseWriter, in *http.Request) {
			h1.Init(r.application, s)

			out.Header().Set("Content-Type", "text/html; charset=utf-8")

			h1.InitRequest(out, in)
			h1.Handle()
			h1.Render()
		}))
	} else if h2, ok := h.(http.Handler); ok {
		chain = alice.New(
			r.loggerMiddleware,
			r.authMiddleware,
		)

		r.Router.Handler(m, p, chain.Then(h2))
	} else if h3, ok := h.(http.HandlerFunc); ok {
		chain = alice.New(
			r.loggerMiddleware,
			r.authMiddleware,
		)

		r.Router.Handler(m, p, chain.ThenFunc(h3))
	} else {
		panic(fmt.Sprintf("Unknown handler type %s %s", m, p))
	}
}
