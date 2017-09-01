package dashboard

import (
	"time"

	"github.com/alexedwards/scs"
	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigHost               = ComponentName + ".host"
	ConfigPort               = ComponentName + ".port"
	ConfigAuthUser           = ComponentName + ".auth-user"
	ConfigAuthPassword       = ComponentName + ".auth-password"
	ConfigSessionCookieName  = ComponentName + ".session.cookie-name"
	ConfigSessionDomain      = ComponentName + ".session.domain"
	ConfigSessionHttpOnly    = ComponentName + ".session.http-only"
	ConfigSessionIdleTimeout = ComponentName + ".session.idle-timeout"
	ConfigSessionLifetime    = ComponentName + ".session.lifetime"
	ConfigSessionPath        = ComponentName + ".session.path"
	ConfigSessionPersist     = ComponentName + ".session.persist"
	ConfigSessionSecure      = ComponentName + ".session.secure"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:     ConfigHost,
			Default: "localhost",
			Usage:   "Frontend host",
			Type:    config.ValueTypeString,
		},
		{
			Key:     ConfigPort,
			Default: 8080,
			Usage:   "Frontend port number",
			Type:    config.ValueTypeInt,
		},
		{
			Key:      ConfigAuthUser,
			Usage:    "User login",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigAuthPassword,
			Usage:    "User password",
			Type:     config.ValueTypeString,
			Editable: true,
			View:     []string{config.ViewPassword},
		},
		{
			Key:      ConfigSessionCookieName,
			Usage:    "The name of the session cookie issued to clients. Note that cookie names should not contain whitespace, commas, semicolons, backslashes or control characters as per RFC6265n",
			Type:     config.ValueTypeString,
			Editable: true,
			Default:  "shadow.session",
		},
		{
			Key:      ConfigSessionDomain,
			Usage:    "Domain attribute on the session cookie",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigSessionHttpOnly,
			Usage:    "HttpOnly attribute on the session cookie",
			Type:     config.ValueTypeBool,
			Editable: true,
			Default:  true,
		},
		{
			Key:      ConfigSessionIdleTimeout,
			Usage:    "Maximum length of time a session can be inactive before it expires",
			Type:     config.ValueTypeDuration,
			Editable: true,
			Default:  0,
		},
		{
			Key:      ConfigSessionLifetime,
			Usage:    "Maximum length of time that a session is valid for before it expires",
			Type:     config.ValueTypeDuration,
			Editable: true,
			Default:  24 * time.Hour,
		},
		{
			Key:      ConfigSessionPath,
			Usage:    "Path attribute on the session cookie",
			Type:     config.ValueTypeString,
			Editable: true,
			Default:  "/",
		},
		{
			Key:      ConfigSessionPersist,
			Usage:    "Persist sets whether the session cookie should be persistent or not",
			Type:     config.ValueTypeBool,
			Editable: true,
			Default:  false,
		},
		{
			Key:      ConfigSessionSecure,
			Usage:    "Secure attribute on the session cookie",
			Type:     config.ValueTypeBool,
			Editable: true,
			Default:  false,
		},
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		config.WatcherForAll:     {c.watchConfig},
		ConfigSessionCookieName:  {c.watchSessionCookieName},
		ConfigSessionDomain:      {c.watchSessionDomain},
		ConfigSessionHttpOnly:    {c.watchSessionHttpOnly},
		ConfigSessionIdleTimeout: {c.watchSessionIdleTimeout},
		ConfigSessionLifetime:    {c.watchSessionLifetime},
		ConfigSessionPath:        {c.watchSessionPath},
		ConfigSessionPersist:     {c.watchSessionPersist},
		ConfigSessionSecure:      {c.watchSessionSecure},
	}
}

func (c *Component) watchConfig(key string, newValue interface{}, oldValue interface{}) {
	c.logger.Infof("Change value for %s with '%v' to '%v'", key, oldValue, newValue)
}

func (c *Component) watchSessionCookieName(_ string, v interface{}, _ interface{}) {
	scs.CookieName = v.(string)
}

func (c *Component) watchSessionDomain(_ string, v interface{}, _ interface{}) {
	c.session.Domain(v.(string))
}

func (c *Component) watchSessionHttpOnly(_ string, v interface{}, _ interface{}) {
	c.session.HttpOnly(v.(bool))
}

func (c *Component) watchSessionIdleTimeout(_ string, v interface{}, _ interface{}) {
	c.session.IdleTimeout(v.(time.Duration))
}

func (c *Component) watchSessionLifetime(_ string, v interface{}, _ interface{}) {
	c.session.Lifetime(v.(time.Duration))
}

func (c *Component) watchSessionPath(_ string, v interface{}, _ interface{}) {
	c.session.Path(v.(string))
}

func (c *Component) watchSessionPersist(_ string, v interface{}, _ interface{}) {
	c.session.Persist(v.(bool))
}

func (c *Component) watchSessionSecure(_ string, v interface{}, _ interface{}) {
	c.session.Secure(v.(bool))
}
