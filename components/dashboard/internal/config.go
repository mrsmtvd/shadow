package internal

import (
	"time"

	"github.com/alexedwards/scs"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariableItem(
			dashboard.ConfigHost,
			config.ValueTypeString,
			"localhost",
			"Frontend host",
			false,
			nil,
			nil),
		config.NewVariableItem(dashboard.ConfigPort,
			config.ValueTypeInt,
			8080,
			"Frontend port number",
			false,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigAuthEnabled,
			config.ValueTypeBool,
			false,
			"Enabled standard auth",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigAuthUser,
			config.ValueTypeString,
			nil,
			"Standard auth user login",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigAuthPassword,
			config.ValueTypeString,
			nil,
			"Standard auth password",
			true,
			[]string{config.ViewPassword},
			nil),
		config.NewVariableItem(
			dashboard.ConfigOAuth2Enabled,
			config.ValueTypeBool,
			nil,
			"Enabled oAuth2",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigOAuth2ID,
			config.ValueTypeString,
			nil,
			"oAuth2 client id",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigOAuth2Secret,
			config.ValueTypeString,
			nil,
			"oAuth2 client secret",
			true,
			[]string{config.ViewPassword},
			nil),
		config.NewVariableItem(
			dashboard.ConfigOAuth2Scopes,
			config.ValueTypeString,
			nil,
			"oAuth2 scopes",
			true,
			[]string{config.ViewTags},
			map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a scope",
			}),
		config.NewVariableItem(
			dashboard.ConfigOAuth2AuthURL,
			config.ValueTypeString,
			nil,
			"oAuth2 endpoint auth URL",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigOAuth2TokenURL,
			config.ValueTypeString,
			nil,
			"oAuth2 endpoint token URL",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigOAuth2ProfileURL,
			config.ValueTypeString,
			nil,
			"oAuth2 endpoint profile URL",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigOAuth2RedirectURL,
			config.ValueTypeString,
			nil,
			"oAuth2 redirect URL",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigSessionCookieName,
			config.ValueTypeString,
			"shadow.session",
			"The name of the session cookie issued to clients. Note that cookie names should not contain whitespace, commas, semicolons, backslashes or control characters as per RFC6265n",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigSessionDomain,
			config.ValueTypeString,
			nil,
			"Domain attribute on the session cookie",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigSessionHttpOnly,
			config.ValueTypeBool,
			true,
			"HttpOnly attribute on the session cookie",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigSessionIdleTimeout,
			config.ValueTypeDuration,
			0,
			"Maximum length of time a session can be inactive before it expires",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigSessionLifetime,
			config.ValueTypeDuration,
			24*time.Hour,
			"Maximum length of time that a session is valid for before it expires",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigSessionPath,
			config.ValueTypeString,
			"/",
			"Path attribute on the session cookie",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigSessionPersist,
			config.ValueTypeBool,
			false,
			"Persist sets whether the session cookie should be persistent or not",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigSessionSecure,
			config.ValueTypeBool,
			false,
			"Secure attribute on the session cookie",
			true,
			nil,
			nil),
		config.NewVariableItem(
			dashboard.ConfigFrontendMinifyEnabled,
			config.ValueTypeBool,
			true,
			"Use minified static files",
			true,
			nil,
			nil),
	}
}

func (c *Component) GetConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher(dashboard.ComponentName, []string{
			dashboard.ConfigAuthEnabled,
			dashboard.ConfigAuthUser,
			dashboard.ConfigAuthPassword,
			dashboard.ConfigOAuth2Enabled,
			dashboard.ConfigOAuth2ID,
			dashboard.ConfigOAuth2Secret,
			dashboard.ConfigOAuth2Scopes,
			dashboard.ConfigOAuth2AuthURL,
			dashboard.ConfigOAuth2TokenURL,
			dashboard.ConfigOAuth2ProfileURL,
			dashboard.ConfigOAuth2RedirectURL,
		}, c.watchAuth),
		config.NewWatcher(dashboard.ComponentName, []string{dashboard.ConfigSessionCookieName}, c.watchSessionCookieName),
		config.NewWatcher(dashboard.ComponentName, []string{dashboard.ConfigSessionDomain}, c.watchSessionDomain),
		config.NewWatcher(dashboard.ComponentName, []string{dashboard.ConfigSessionHttpOnly}, c.watchSessionHttpOnly),
		config.NewWatcher(dashboard.ComponentName, []string{dashboard.ConfigSessionIdleTimeout}, c.watchSessionIdleTimeout),
		config.NewWatcher(dashboard.ComponentName, []string{dashboard.ConfigSessionLifetime}, c.watchSessionLifetime),
		config.NewWatcher(dashboard.ComponentName, []string{dashboard.ConfigSessionPath}, c.watchSessionPath),
		config.NewWatcher(dashboard.ComponentName, []string{dashboard.ConfigSessionPersist}, c.watchSessionPersist),
		config.NewWatcher(dashboard.ComponentName, []string{dashboard.ConfigSessionSecure}, c.watchSessionSecure),
	}
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
