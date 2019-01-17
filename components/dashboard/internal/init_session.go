package internal

import (
	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/stores/memstore"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) initSession() {
	store := memstore.New(0)
	c.session = scs.NewManager(store)
	c.session.Name(c.config.String(dashboard.ConfigSessionCookieName))

	c.session.Domain(c.config.String(dashboard.ConfigSessionDomain))
	c.session.HttpOnly(c.config.Bool(dashboard.ConfigSessionHttpOnly))
	c.session.IdleTimeout(c.config.Duration(dashboard.ConfigSessionIdleTimeout))
	c.session.Lifetime(c.config.Duration(dashboard.ConfigSessionLifetime))
	c.session.Path(c.config.String(dashboard.ConfigSessionPath))
	c.session.Persist(c.config.Bool(dashboard.ConfigSessionPersist))
	c.session.Secure(c.config.Bool(dashboard.ConfigSessionSecure))
}
