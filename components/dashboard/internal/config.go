package internal

import (
	"time"

	"github.com/alexedwards/scs"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			dashboard.ConfigHost,
			config.ValueTypeString,
			"localhost",
			"Frontend host",
			false,
			nil,
			nil),
		config.NewVariable(dashboard.ConfigPort,
			config.ValueTypeInt,
			8080,
			"Frontend port number",
			false,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigAuthEnabled,
			config.ValueTypeBool,
			false,
			"Enabled standard auth",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigAuthUser,
			config.ValueTypeString,
			nil,
			"Standard auth user login",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigAuthPassword,
			config.ValueTypeString,
			nil,
			"Standard auth password",
			true,
			[]string{config.ViewPassword},
			nil),
		config.NewVariable(
			dashboard.ConfigOAuth2EmailsAllowed,
			config.ValueTypeString,
			nil,
			"OAuth emails allowed",
			true,
			[]string{config.ViewTags},
			map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a email",
			}),
		config.NewVariable(
			dashboard.ConfigOAuth2GithubEnabled,
			config.ValueTypeBool,
			nil,
			"Enabled Github provider",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigOAuth2GithubID,
			config.ValueTypeString,
			nil,
			"Github client id",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigOAuth2GithubSecret,
			config.ValueTypeString,
			nil,
			"Github client secret",
			true,
			[]string{config.ViewPassword},
			nil),
		config.NewVariable(
			dashboard.ConfigOAuth2GithubScopes,
			config.ValueTypeString,
			nil,
			"Github scopes",
			true,
			[]string{config.ViewTags},
			map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a scope",
			}),
		config.NewVariable(
			dashboard.ConfigOAuth2GithubRedirectURL,
			config.ValueTypeString,
			nil,
			"Github redirect URL",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigOAuth2GitlabEnabled,
			config.ValueTypeBool,
			nil,
			"Enabled Gitlab provider",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigOAuth2GitlabID,
			config.ValueTypeString,
			nil,
			"Gitlab client id",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigOAuth2GitlabSecret,
			config.ValueTypeString,
			nil,
			"Gitlab client secret",
			true,
			[]string{config.ViewPassword},
			nil),
		config.NewVariable(
			dashboard.ConfigOAuth2GitlabScopes,
			config.ValueTypeString,
			"read_user",
			"Gitlab scopes",
			true,
			[]string{config.ViewTags},
			map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a scope",
			}),
		config.NewVariable(
			dashboard.ConfigOAuth2GitlabAuthURL,
			config.ValueTypeString,
			nil,
			"Gitlab endpoint auth URL",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigOAuth2GitlabTokenURL,
			config.ValueTypeString,
			nil,
			"Gitlab endpoint token URL",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigOAuth2GitlabProfileURL,
			config.ValueTypeString,
			nil,
			"Gitlab endpoint profile URL",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigOAuth2GitlabRedirectURL,
			config.ValueTypeString,
			nil,
			"Gitlab redirect URL",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigSessionCookieName,
			config.ValueTypeString,
			"shadow.session",
			"The name of the session cookie issued to clients. Note that cookie names should not contain whitespace, commas, semicolons, backslashes or control characters as per RFC6265n",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigSessionDomain,
			config.ValueTypeString,
			nil,
			"Domain attribute on the session cookie",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigSessionHttpOnly,
			config.ValueTypeBool,
			true,
			"HttpOnly attribute on the session cookie",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigSessionIdleTimeout,
			config.ValueTypeDuration,
			0,
			"Maximum length of time a session can be inactive before it expires",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigSessionLifetime,
			config.ValueTypeDuration,
			24*time.Hour,
			"Maximum length of time that a session is valid for before it expires",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigSessionPath,
			config.ValueTypeString,
			"/",
			"Path attribute on the session cookie",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigSessionPersist,
			config.ValueTypeBool,
			false,
			"Persist sets whether the session cookie should be persistent or not",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigSessionSecure,
			config.ValueTypeBool,
			false,
			"Secure attribute on the session cookie",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigFrontendMinifyEnabled,
			config.ValueTypeBool,
			true,
			"Use minified static files",
			true,
			nil,
			nil),
		config.NewVariable(
			dashboard.ConfigStartURL,
			config.ValueTypeString,
			"/"+c.GetName(),
			"Start URL",
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
			dashboard.ConfigOAuth2GithubEnabled,
			dashboard.ConfigOAuth2GitlabEnabled,
			dashboard.ConfigOAuth2GitlabID,
			dashboard.ConfigOAuth2GitlabSecret,
			dashboard.ConfigOAuth2GitlabScopes,
			dashboard.ConfigOAuth2GitlabAuthURL,
			dashboard.ConfigOAuth2GitlabTokenURL,
			dashboard.ConfigOAuth2GitlabProfileURL,
			dashboard.ConfigOAuth2GitlabRedirectURL,
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
