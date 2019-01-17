package internal

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/auth"
	"github.com/kihamo/shadow/components/dashboard/auth/providers/password"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/gplus"
)

func (c *Component) initAuth() (err error) {
	auth.ClearProviders()
	providers := make([]goth.Provider, 0)

	var baseURL *url.URL
	baseURLFromConfig := c.config.String(dashboard.ConfigOAuth2BaseURL)
	if baseURLFromConfig != "" {
		if baseURL, err = url.Parse(baseURLFromConfig); err != nil {
			return err
		}
	}

	if baseURL == nil {
		return errors.New("base path for auth callbacks is empty")
	}

	baseURL.Path = strings.Trim(baseURL.Path, "/")
	pathCallbackTpl := "%s/" + strings.Trim(dashboard.AuthPath, "/") + "/%s/callback"

	if c.config.Bool(dashboard.ConfigAuthEnabled) {
		passwordRedirectURL := new(url.URL)
		*passwordRedirectURL = *baseURL
		passwordRedirectURL.Path = fmt.Sprintf(pathCallbackTpl, passwordRedirectURL.Path, "github")

		providers = append(providers, password.New(
			passwordRedirectURL.String(),
			map[string]string{
				c.config.String(dashboard.ConfigAuthUser): c.config.String(dashboard.ConfigAuthPassword),
			},
		))
	}

	if c.config.Bool(dashboard.ConfigOAuth2GithubEnabled) {
		githubRedirectURL := new(url.URL)
		*githubRedirectURL = *baseURL
		githubRedirectURL.Path = fmt.Sprintf(pathCallbackTpl, githubRedirectURL.Path, "github")

		providers = append(providers, github.NewCustomisedURL(
			c.config.String(dashboard.ConfigOAuth2GithubID),
			c.config.String(dashboard.ConfigOAuth2GithubSecret),
			githubRedirectURL.String(),
			github.AuthURL,
			github.TokenURL,
			github.ProfileURL,
			github.EmailURL,
			strings.Split(c.config.String(dashboard.ConfigOAuth2GithubScopes), ",")...,
		))
	}

	if c.config.Bool(dashboard.ConfigOAuth2GitlabEnabled) {
		gitlabRedirectURL := new(url.URL)
		*gitlabRedirectURL = *baseURL
		gitlabRedirectURL.Path = fmt.Sprintf(pathCallbackTpl, gitlabRedirectURL.Path, "gitlab")

		providers = append(providers, gitlab.NewCustomisedURL(
			c.config.String(dashboard.ConfigOAuth2GitlabID),
			c.config.String(dashboard.ConfigOAuth2GitlabSecret),
			gitlabRedirectURL.String(),
			c.config.String(dashboard.ConfigOAuth2GitlabAuthURL),
			c.config.String(dashboard.ConfigOAuth2GitlabTokenURL),
			c.config.String(dashboard.ConfigOAuth2GitlabProfileURL),
			strings.Split(c.config.String(dashboard.ConfigOAuth2GitlabScopes), ",")...,
		))
	}

	if c.config.Bool(dashboard.ConfigOAuth2GplusEnabled) {
		gplusRedirectURL := new(url.URL)
		*gplusRedirectURL = *baseURL
		gplusRedirectURL.Path = fmt.Sprintf(pathCallbackTpl, gplusRedirectURL.Path, "gplus")

		providers = append(providers, gplus.New(
			c.config.String(dashboard.ConfigOAuth2GplusID),
			c.config.String(dashboard.ConfigOAuth2GplusSecret),
			gplusRedirectURL.String(),
			strings.Split(c.config.String(dashboard.ConfigOAuth2GplusScopes), ",")...,
		))
	}

	auth.UseProviders(providers...)

	return nil
}
