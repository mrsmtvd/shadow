package internal

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) initSession() {
	c.sessionManager = scs.New()

	c.sessionManager.IdleTimeout = c.config.Duration(dashboard.ConfigSessionIdleTimeout)
	c.sessionManager.Lifetime = c.config.Duration(dashboard.ConfigSessionLifetime)

	if d := c.config.Duration(dashboard.ConfigSessionCleanupInterval); d > 0 {
		c.sessionManager.Store = memstore.NewWithCleanupInterval(d)
	} else {
		c.sessionManager.Store = memstore.New()
	}

	c.sessionManager.Cookie.Name = c.config.String(dashboard.ConfigSessionCookieName)
	c.sessionManager.Cookie.Domain = c.config.String(dashboard.ConfigSessionDomain)
	c.sessionManager.Cookie.HttpOnly = c.config.Bool(dashboard.ConfigSessionHTTPOnly)
	c.sessionManager.Cookie.Path = c.config.String(dashboard.ConfigSessionPath)
	c.sessionManager.Cookie.Persist = c.config.Bool(dashboard.ConfigSessionPersist)
	c.sessionManager.Cookie.SameSite = http.SameSite(c.config.Int(dashboard.ConfigSessionSameSite))
	c.sessionManager.Cookie.Secure = c.config.Bool(dashboard.ConfigSessionSecure)

	c.sessionManager.ErrorFunc = func(w http.ResponseWriter, r *http.Request, err error) {
		c.logger.Error(err.Error())
		c.router.InternalErrorServeHTTP(w, r, err)
	}
}
