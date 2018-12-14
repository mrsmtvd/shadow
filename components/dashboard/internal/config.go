package internal

import (
	"time"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(dashboard.ConfigHost, config.ValueTypeString).
			WithUsage("Host").
			WithGroup("Listen").
			WithDefault("localhost"),
		config.NewVariable(dashboard.ConfigPort, config.ValueTypeInt).
			WithUsage("Port number").
			WithGroup("Listen").
			WithDefault(8080),
		config.NewVariable(dashboard.ConfigAuthEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("Authorization basic").
			WithEditable(true),
		config.NewVariable(dashboard.ConfigAuthUser, config.ValueTypeString).
			WithUsage("User login").
			WithGroup("Authorization basic").
			WithEditable(true),
		config.NewVariable(dashboard.ConfigAuthPassword, config.ValueTypeString).
			WithUsage("User password").
			WithGroup("Authorization basic").
			WithEditable(true).
			WithView([]string{config.ViewPassword}),
		config.NewVariable(dashboard.ConfigOAuth2EmailsAllowed, config.ValueTypeString).
			WithUsage("Emails allowed").
			WithGroup("Authorization OAuth").
			WithEditable(true).
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add a email"}),
		config.NewVariable(dashboard.ConfigOAuth2DomainsAllowed, config.ValueTypeString).
			WithUsage("Domains of emails allowed").
			WithGroup("Authorization OAuth").
			WithEditable(true).
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add a domain"}),
		config.NewVariable(dashboard.ConfigOAuth2BaseURL, config.ValueTypeString).
			WithUsage("Base URL for redirect").
			WithGroup("Authorization OAuth").
			WithEditable(true).
			WithDefault("http://localhost/"),
		config.NewVariable(dashboard.ConfigOAuth2GithubEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("Authorization OAuth Github provider").
			WithEditable(true),
		config.NewVariable(dashboard.ConfigOAuth2GithubID, config.ValueTypeString).
			WithUsage("Client ID").
			WithGroup("Authorization OAuth Github provider").
			WithEditable(true),
		config.NewVariable(dashboard.ConfigOAuth2GithubSecret, config.ValueTypeString).
			WithUsage("Client secret").
			WithGroup("Authorization OAuth Github provider").
			WithEditable(true).
			WithView([]string{config.ViewPassword}),
		config.NewVariable(dashboard.ConfigOAuth2GithubScopes, config.ValueTypeString).
			WithUsage("Scopes").
			WithGroup("Authorization OAuth Github provider").
			WithEditable(true).
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add a scope"}),
		config.NewVariable(dashboard.ConfigOAuth2GitlabEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("Authorization OAuth Gitlab provider").
			WithEditable(true),
		config.NewVariable(dashboard.ConfigOAuth2GitlabID, config.ValueTypeString).
			WithUsage("Client ID").
			WithGroup("Authorization OAuth Gitlab provider").
			WithEditable(true),
		config.NewVariable(dashboard.ConfigOAuth2GitlabSecret, config.ValueTypeString).
			WithUsage("Client secret").
			WithGroup("Authorization OAuth Gitlab provider").
			WithEditable(true).
			WithView([]string{config.ViewPassword}),
		config.NewVariable(dashboard.ConfigOAuth2GitlabScopes, config.ValueTypeString).
			WithUsage("Scopes").
			WithGroup("Authorization OAuth Gitlab provider").
			WithEditable(true).
			WithDefault("read_user").
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add a scope"}),
		config.NewVariable(dashboard.ConfigOAuth2GitlabAuthURL, config.ValueTypeString).
			WithUsage("Endpoint auth URL").
			WithGroup("Authorization OAuth Gitlab provider").
			WithEditable(true),
		config.NewVariable(dashboard.ConfigOAuth2GitlabTokenURL, config.ValueTypeString).
			WithUsage("Endpoint token URL").
			WithGroup("Authorization OAuth Gitlab provider").
			WithEditable(true),
		config.NewVariable(dashboard.ConfigOAuth2GitlabProfileURL, config.ValueTypeString).
			WithUsage("Endpoint profile URL").
			WithGroup("Authorization OAuth Gitlab provider").
			WithEditable(true),
		config.NewVariable(dashboard.ConfigOAuth2GplusEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("Authorization OAuth Google+ provider").
			WithEditable(true),
		config.NewVariable(dashboard.ConfigOAuth2GplusID, config.ValueTypeString).
			WithUsage("Client ID").
			WithGroup("Authorization OAuth Google+ provider").
			WithEditable(true),
		config.NewVariable(dashboard.ConfigOAuth2GplusSecret, config.ValueTypeString).
			WithUsage("Client secret").
			WithGroup("Authorization OAuth Google+ provider").
			WithEditable(true).
			WithView([]string{config.ViewPassword}),
		config.NewVariable(dashboard.ConfigOAuth2GplusScopes, config.ValueTypeString).
			WithUsage("Scopes").
			WithGroup("Authorization OAuth Google+ provider").
			WithEditable(true).
			WithDefault("profile,email,openid").
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add a scope"}),
		config.NewVariable(dashboard.ConfigSessionCookieName, config.ValueTypeString).
			WithUsage("The name of the session cookie issued to clients. Note that cookie names should not contain whitespace, commas, semicolons, backslashes or control characters as per RFC6265n").
			WithGroup("User session").
			WithEditable(true).
			WithDefault("shadow.session"),
		config.NewVariable(dashboard.ConfigSessionDomain, config.ValueTypeString).
			WithUsage("Domain attribute on the session cookie").
			WithGroup("User session").
			WithEditable(true),
		config.NewVariable(dashboard.ConfigSessionHttpOnly, config.ValueTypeBool).
			WithUsage("HttpOnly attribute on the session cookie").
			WithGroup("User session").
			WithEditable(true).
			WithDefault(true),
		config.NewVariable(dashboard.ConfigSessionIdleTimeout, config.ValueTypeDuration).
			WithUsage("Maximum length of time a session can be inactive before it expires").
			WithGroup("User session").
			WithEditable(true).
			WithDefault(0),
		config.NewVariable(dashboard.ConfigSessionLifetime, config.ValueTypeDuration).
			WithUsage("Maximum length of time that a session is valid for before it expires").
			WithGroup("User session").
			WithEditable(true).
			WithDefault(24 * time.Hour),
		config.NewVariable(dashboard.ConfigSessionPath, config.ValueTypeString).
			WithUsage("Path attribute on the session cookie").
			WithGroup("User session").
			WithEditable(true).
			WithDefault("/"),
		config.NewVariable(dashboard.ConfigSessionPersist, config.ValueTypeBool).
			WithUsage("Persist sets whether the session cookie should be persistent or not").
			WithGroup("User session").
			WithEditable(true),
		config.NewVariable(dashboard.ConfigSessionSecure, config.ValueTypeBool).
			WithUsage("Secure attribute on the session cookie").
			WithGroup("User session").
			WithEditable(true),
		config.NewVariable(dashboard.ConfigFrontendMinifyEnabled, config.ValueTypeBool).
			WithUsage("Use minified static files").
			WithGroup("Develop mode").
			WithEditable(true).
			WithDefault(true),
		config.NewVariable(dashboard.ConfigPanicHandlerCallerSkip, config.ValueTypeInt64).
			WithUsage("Skip number of callers in panic handler").
			WithGroup("Develop mode").
			WithEditable(true).
			WithDefault(DefaultCallerSkip),
		config.NewVariable(dashboard.ConfigStartURL, config.ValueTypeString).
			WithUsage("Start URL").
			WithEditable(true).
			WithDefault("/" + c.Name()),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{
			dashboard.ConfigAuthEnabled,
			dashboard.ConfigAuthUser,
			dashboard.ConfigAuthPassword,
			dashboard.ConfigOAuth2GithubEnabled,
			dashboard.ConfigOAuth2GithubID,
			dashboard.ConfigOAuth2GithubSecret,
			dashboard.ConfigOAuth2GithubScopes,
			dashboard.ConfigOAuth2GitlabEnabled,
			dashboard.ConfigOAuth2GitlabID,
			dashboard.ConfigOAuth2GitlabSecret,
			dashboard.ConfigOAuth2GitlabScopes,
			dashboard.ConfigOAuth2GitlabAuthURL,
			dashboard.ConfigOAuth2GitlabTokenURL,
			dashboard.ConfigOAuth2GitlabProfileURL,
			dashboard.ConfigOAuth2GplusEnabled,
			dashboard.ConfigOAuth2GplusID,
			dashboard.ConfigOAuth2GplusSecret,
			dashboard.ConfigOAuth2GplusScopes,
		}, c.watchAuth),
		config.NewWatcher([]string{dashboard.ConfigSessionCookieName}, c.watchSessionCookieName),
		config.NewWatcher([]string{dashboard.ConfigSessionDomain}, c.watchSessionDomain),
		config.NewWatcher([]string{dashboard.ConfigSessionHttpOnly}, c.watchSessionHttpOnly),
		config.NewWatcher([]string{dashboard.ConfigSessionIdleTimeout}, c.watchSessionIdleTimeout),
		config.NewWatcher([]string{dashboard.ConfigSessionLifetime}, c.watchSessionLifetime),
		config.NewWatcher([]string{dashboard.ConfigSessionPath}, c.watchSessionPath),
		config.NewWatcher([]string{dashboard.ConfigSessionPersist}, c.watchSessionPersist),
		config.NewWatcher([]string{dashboard.ConfigSessionSecure}, c.watchSessionSecure),
		config.NewWatcher([]string{dashboard.ConfigPanicHandlerCallerSkip}, c.watchPanicHandlerCallerSkip),
	}
}

func (c *Component) watchAuth(_ string, _ interface{}, _ interface{}) {
	_ = c.initAuth()
}

func (c *Component) watchSessionCookieName(_ string, v interface{}, _ interface{}) {
	c.session.Name(v.(string))
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

func (c *Component) watchPanicHandlerCallerSkip(_ string, v interface{}, _ interface{}) {
	c.router.SetPanicHandlerCallerSkip(int(v.(int64)))
}
