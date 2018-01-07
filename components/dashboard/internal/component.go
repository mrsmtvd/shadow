package internal

import (
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/stores/memstore"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/auth"
	"github.com/kihamo/shadow/components/dashboard/auth/providers/password"
	"github.com/kihamo/shadow/components/logger"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logger.Logger
	renderer    *Renderer
	session     *scs.Manager
	routes      []dashboard.Route
}

func (c *Component) GetName() string {
	return dashboard.ComponentName
}

func (c *Component) GetVersion() string {
	return dashboard.ComponentVersion
}

func (c *Component) GetDependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: logger.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) (err error) {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) error {
	c.logger = logger.NewOrNop(c.GetName(), c.application)

	if err := c.loadTemplates(); err != nil {
		return err
	}

	if err := c.loadMenu(); err != nil {
		return err
	}

	c.initSession()
	c.initAuth()

	mux, err := c.getServeMux()
	if err != nil {
		return err
	}

	go func() {
		defer wg.Done()

		addr := net.JoinHostPort(c.config.GetString(dashboard.ConfigHost), c.config.GetString(dashboard.ConfigPort))
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			c.logger.Fatalf("Failed to listen [%d]: %s\n", os.Getpid(), err.Error())
		}

		c.logger.Info("Running service", map[string]interface{}{
			"addr": addr,
			"pid":  os.Getpid(),
		})

		if err := http.Serve(lis, mux); err != nil {
			c.logger.Fatalf("Failed to serve [%d]: %s\n", os.Getpid(), err.Error())
		}
	}()

	return nil
}

func (c *Component) initAuth() {
	auth.ClearProviders()
	providers := []goth.Provider{}

	if c.config.GetBool(dashboard.ConfigAuthEnabled) {
		providers = append(providers, password.New(
			dashboard.AuthPath+"/password/callback",
			map[string]string{
				c.config.GetString(dashboard.ConfigAuthUser): c.config.GetString(dashboard.ConfigAuthPassword),
			},
		))
	}

	if c.config.GetBool(dashboard.ConfigOAuth2GithubEnabled) {
		providers = append(providers, github.NewCustomisedURL(
			c.config.GetString(dashboard.ConfigOAuth2GithubID),
			c.config.GetString(dashboard.ConfigOAuth2GithubSecret),
			c.config.GetString(dashboard.ConfigOAuth2GithubRedirectURL),
			github.AuthURL,
			github.TokenURL,
			github.ProfileURL,
			github.EmailURL,
			strings.Split(c.config.GetString(dashboard.ConfigOAuth2GithubScopes), ",")...,
		))
	}

	if c.config.GetBool(dashboard.ConfigOAuth2GitlabEnabled) {
		providers = append(providers, gitlab.NewCustomisedURL(
			c.config.GetString(dashboard.ConfigOAuth2GitlabID),
			c.config.GetString(dashboard.ConfigOAuth2GitlabSecret),
			c.config.GetString(dashboard.ConfigOAuth2GitlabRedirectURL),
			c.config.GetString(dashboard.ConfigOAuth2GitlabAuthURL),
			c.config.GetString(dashboard.ConfigOAuth2GitlabTokenURL),
			c.config.GetString(dashboard.ConfigOAuth2GitlabProfileURL),
			strings.Split(c.config.GetString(dashboard.ConfigOAuth2GitlabScopes), ",")...,
		))
	}

	auth.UseProviders(providers...)
}

func (c *Component) initSession() {
	scs.CookieName = c.config.GetString(dashboard.ConfigSessionCookieName)

	store := memstore.New(0)
	c.session = scs.NewManager(store)

	c.session.Domain(c.config.GetString(dashboard.ConfigSessionDomain))
	c.session.HttpOnly(c.config.GetBool(dashboard.ConfigSessionHttpOnly))
	c.session.IdleTimeout(c.config.GetDuration(dashboard.ConfigSessionIdleTimeout))
	c.session.Lifetime(c.config.GetDuration(dashboard.ConfigSessionLifetime))
	c.session.Path(c.config.GetString(dashboard.ConfigSessionPath))
	c.session.Persist(c.config.GetBool(dashboard.ConfigSessionPersist))
	c.session.Secure(c.config.GetBool(dashboard.ConfigSessionSecure))
}
