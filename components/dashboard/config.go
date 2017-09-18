package dashboard

import (
	"time"

	"github.com/alexedwards/scs"
	"github.com/kihamo/shadow/components/config"
)

const (
	ConfigHost               = ComponentName + ".host"
	ConfigPort               = ComponentName + ".port"
	ConfigAuthEnabled        = ComponentName + ".auth.enabled"
	ConfigAuthUser           = ComponentName + ".auth.user"
	ConfigAuthPassword       = ComponentName + ".auth.password"
	ConfigOAuth2Enabled      = ComponentName + ".oauth2.enabled"
	ConfigOAuth2ID           = ComponentName + ".oauth2.id"
	ConfigOAuth2Secret       = ComponentName + ".oauth2.secret"
	ConfigOAuth2Scopes       = ComponentName + ".oauth2.scopes"
	ConfigOAuth2AuthURL      = ComponentName + ".oauth2.auth-url"
	ConfigOAuth2TokenURL     = ComponentName + ".oauth2.token-url"
	ConfigOAuth2ProfileURL   = ComponentName + ".oauth2.profile-url"
	ConfigOAuth2RedirectURL  = ComponentName + ".oauth2.redirect-url"
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
			Key:      ConfigAuthEnabled,
			Usage:    "Enabled standard auth",
			Type:     config.ValueTypeBool,
			Editable: true,
			Default:  true,
		},
		{
			Key:      ConfigAuthUser,
			Usage:    "Standard auth user login",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigAuthPassword,
			Usage:    "Standard auth password",
			Type:     config.ValueTypeString,
			Editable: true,
			View:     []string{config.ViewPassword},
		},
		{
			Key:      ConfigOAuth2Enabled,
			Usage:    "Enabled oAuth2",
			Type:     config.ValueTypeBool,
			Editable: true,
		},
		{
			Key:      ConfigOAuth2ID,
			Usage:    "oAuth2 client id",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigOAuth2Secret,
			Usage:    "oAuth2 client secret",
			Type:     config.ValueTypeString,
			Editable: true,
			View:     []string{config.ViewPassword},
		},
		{
			Key:      ConfigOAuth2Scopes,
			Usage:    "oAuth2 scopes",
			Type:     config.ValueTypeString,
			Editable: true,
			View:     []string{config.ViewTags},
			ViewOptions: map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a scope",
			},
		},
		{
			Key:      ConfigOAuth2AuthURL,
			Usage:    "oAuth2 endpoint auth URL",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigOAuth2TokenURL,
			Usage:    "oAuth2 endpoint token URL",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigOAuth2ProfileURL,
			Usage:    "oAuth2 endpoint profile URL",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigOAuth2RedirectURL,
			Usage:    "oAuth2 redirect URL",
			Type:     config.ValueTypeString,
			Editable: true,
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
		ConfigAuthEnabled:        {c.watchAuth},
		ConfigAuthUser:           {c.watchAuth},
		ConfigAuthPassword:       {c.watchAuth},
		ConfigOAuth2Enabled:      {c.watchAuth},
		ConfigOAuth2ID:           {c.watchAuth},
		ConfigOAuth2Secret:       {c.watchAuth},
		ConfigOAuth2Scopes:       {c.watchAuth},
		ConfigOAuth2AuthURL:      {c.watchAuth},
		ConfigOAuth2TokenURL:     {c.watchAuth},
		ConfigOAuth2ProfileURL:   {c.watchAuth},
		ConfigOAuth2RedirectURL:  {c.watchAuth},
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

func (c *Component) watchAuth(_ string, _ interface{}, _ interface{}) {
	c.initAuth()
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
