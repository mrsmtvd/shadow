package dashboard

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
	"github.com/kihamo/shadow/components/dashboard/auth"
	"github.com/kihamo/shadow/components/dashboard/auth/providers/password"
	"github.com/kihamo/shadow/components/logger"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gitlab"
)

const (
	ComponentName = "dashboard"
)

type Component struct {
	application shadow.Application
	config      *config.Component
	logger      logger.Logger
	renderer    *Renderer
	session     *scs.Manager
}

func (c *Component) GetName() string {
	return ComponentName
}

func (c *Component) GetVersion() string {
	return ComponentVersion
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
	c.config = a.GetComponent(config.ComponentName).(*config.Component)

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

		addr := net.JoinHostPort(c.config.GetString(ConfigHost), c.config.GetString(ConfigPort))
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

	if c.config.GetBool(ConfigAuthEnabled) {
		providers = append(providers, password.New(
			AuthPath+"/password/callback",
			map[string]string{
				c.config.GetString(ConfigAuthUser): c.config.GetString(ConfigAuthPassword),
			},
		))
	}

	if c.config.GetBool(ConfigOAuth2Enabled) {
		providers = append(providers, gitlab.NewCustomisedURL(
			c.config.GetString(ConfigOAuth2ID),
			c.config.GetString(ConfigOAuth2Secret),
			c.config.GetString(ConfigOAuth2RedirectURL),
			c.config.GetString(ConfigOAuth2AuthURL),
			c.config.GetString(ConfigOAuth2TokenURL),
			c.config.GetString(ConfigOAuth2ProfileURL),
			strings.Split(c.config.GetString(ConfigOAuth2Scopes), ",")...,
		))
	}

	auth.UseProviders(providers...)
}

func (c *Component) initSession() {
	scs.CookieName = c.config.GetString(ConfigSessionCookieName)

	store := memstore.New(0)
	c.session = scs.NewManager(store)

	c.session.Domain(c.config.GetString(ConfigSessionDomain))
	c.session.HttpOnly(c.config.GetBool(ConfigSessionHttpOnly))
	c.session.IdleTimeout(c.config.GetDuration(ConfigSessionIdleTimeout))
	c.session.Lifetime(c.config.GetDuration(ConfigSessionLifetime))
	c.session.Path(c.config.GetString(ConfigSessionPath))
	c.session.Persist(c.config.GetBool(ConfigSessionPersist))
	c.session.Secure(c.config.GetBool(ConfigSessionSecure))
}
