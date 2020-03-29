package internal

import (
	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/stores/memstore"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) initSession() {
	store := memstore.New(0)
	c.sessionManager = scs.NewManager(store)
	c.sessionManager.Name(c.config.String(dashboard.ConfigSessionCookieName))

	c.sessionManager.Domain(c.config.String(dashboard.ConfigSessionDomain))
	c.sessionManager.HttpOnly(c.config.Bool(dashboard.ConfigSessionHTTPOnly))
	c.sessionManager.IdleTimeout(c.config.Duration(dashboard.ConfigSessionIdleTimeout))
	c.sessionManager.Lifetime(c.config.Duration(dashboard.ConfigSessionLifetime))
	c.sessionManager.Path(c.config.String(dashboard.ConfigSessionPath))
	c.sessionManager.Persist(c.config.Bool(dashboard.ConfigSessionPersist))
	c.sessionManager.Secure(c.config.Bool(dashboard.ConfigSessionSecure))
}
