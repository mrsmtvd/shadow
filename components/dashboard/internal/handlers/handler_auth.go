package handlers

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/dashboard/auth"
	"github.com/markbates/goth"
)

var stateRand = rand.New(rand.NewSource(time.Now().UnixNano()))

type ProvidersView []ProviderView

type ProviderView struct {
	Name         string
	CallbackURL  string
	NeedRedirect bool
}

type AuthHandler struct {
	dashboard.Handler

	IsCallback bool
	Config     config.Component
}

func (h *AuthHandler) buildProvidersView(r *dashboard.Request) ProvidersView {
	providers := make(ProvidersView, 0, len(auth.GetProviders()))

	for _, p := range auth.GetProviders() {
		provider := ProviderView{
			Name:         p.Name(),
			CallbackURL:  "/dashboard/auth/" + p.Name(),
			NeedRedirect: true,
		}

		if p.Name() == "password" {
			callback, err := h.redirectToExternal(r, p)
			if err != nil {
				r.Logger().Errorf("Error get redirect url for %s provider", p.Name(), map[string]interface{}{
					"error": err.Error(),
				})

				continue
			}

			provider.NeedRedirect = false
			provider.CallbackURL = callback
		}

		providers = append(providers, provider)
	}

	return providers
}

func (h *AuthHandler) renderForm(r *dashboard.Request, err error) {
	providers := h.buildProvidersView(r)

	data := map[string]interface{}{
		"error":             err,
		"providers":         providers,
		"hasMultiProviders": false,
	}

	last := false
	for i, p := range providers {
		if i < 1 {
			continue
		} else if p.NeedRedirect != last {
			data["hasMultiProviders"] = true
			break
		}

		last = p.NeedRedirect
	}

	h.RenderLayout(r.Context(), dashboard.ComponentName, "auth", "blank", data)
}

func (h *AuthHandler) getRedirectToLastURL(r *dashboard.Request) string {
	redirectURL, err := r.Session().GetString(dashboard.SessionLastURL)
	if err == nil && redirectURL != "" && !strings.HasPrefix(redirectURL, dashboard.AuthPath) {
		return redirectURL
	}

	return "/"
}

func (h *AuthHandler) auth(r *dashboard.Request, provider goth.Provider) error {
	session := r.Session()
	sessionKey := dashboard.AuthSessionName()

	exists, err := session.Exists(sessionKey)
	if !exists {
		return errors.New("OAuth session not exists")
	}

	value, err := session.GetString(sessionKey)
	if err != nil {
		return err
	}

	providerSession, err := provider.UnmarshalSession(value)
	if err != nil {
		return err
	}

	rawAuthURL, err := providerSession.GetAuthURL()
	if err != nil {
		return err
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return err
	}

	r.Original().ParseForm()
	originalState := authURL.Query().Get("state")

	if originalState != "" && (originalState != r.Original().Form.Get("state")) {
		return errors.New("State token mismatch")
	}

	if providerUser, err := provider.FetchUser(providerSession); err == nil {
		if err = session.PutObject(dashboard.SessionUser, auth.NewUser(providerUser)); err != nil {
			return err
		}

		return nil
	}

	if _, err = providerSession.Authorize(provider, r.Original().Form); err != nil {
		return err
	}

	providerUser, err := provider.FetchUser(providerSession)
	if err != nil {
		return err
	}

	if provider.Name() != "password" {
		emailsConfig := h.Config.GetString(dashboard.ConfigOAuth2EmailsAllowed)
		if emailsConfig != "" {
			emails := strings.Split(emailsConfig, ",")

			if len(emails) > 0 {
				var valid bool

				if providerUser.Email != "" {
					for _, email := range emails {
						if email == providerUser.Email {
							valid = true
							break
						}
					}
				}

				if !valid {
					return errors.New("Email not allowed")
				}
			}
		}
	}

	if err = session.PutString(sessionKey, providerSession.Marshal()); err != nil {
		return err
	}

	r.Logger().Debugf("Auth user %s is success", providerUser.Name, map[string]interface{}{
		"auth.provider":            provider.Name(),
		"auth.user-id":             providerUser.UserID,
		"auth.email":               providerUser.Email,
		"auth.access-token":        providerUser.AccessToken,
		"auth.access-token-secret": providerUser.AccessTokenSecret,
		"auth.refresh-token":       providerUser.RefreshToken,
		"auth.expires":             providerUser.ExpiresAt,
	})

	if err = session.PutObject(dashboard.SessionUser, auth.NewUser(providerUser)); err != nil {
		return err
	}

	return nil
}

func (h *AuthHandler) redirectToExternal(r *dashboard.Request, provider goth.Provider) (string, error) {
	state := r.URL().Query().Get("state")
	if len(state) == 0 {
		nonceBytes := make([]byte, 64)
		for i := 0; i < 64; i++ {
			nonceBytes[i] = byte(stateRand.Int63() % 256)
		}

		state = base64.URLEncoding.EncodeToString(nonceBytes)
	}

	providerSession, err := provider.BeginAuth(state)
	if err != nil {
		return "", err
	}

	externalUrl, err := providerSession.GetAuthURL()
	if err != nil {
		return "", err
	}

	if externalUrl == "" {
		return "", errors.New("External url for redirect is empty")
	}

	if err = r.Session().PutString(dashboard.AuthSessionName(), providerSession.Marshal()); err != nil {
		return "", err
	}

	return externalUrl, nil
}

func (h *AuthHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if len(auth.GetProviders()) == 0 || r.User().IsAuthorized() {
		h.Redirect(h.getRedirectToLastURL(r), http.StatusSeeOther, w, r)
		return
	}

	providerName := r.URL().Query().Get(":provider")
	if providerName == "" {
		if !h.IsCallback {
			h.renderForm(r, nil)
		} else {
			h.NotFound(w, r)
		}

		return
	}

	provider, err := auth.GetProvider(providerName)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	if !h.IsCallback {
		externalUrl, err := h.redirectToExternal(r, provider)
		if err != nil {
			h.renderForm(r, err)
			return
		}

		r.Logger().Debugf("OAuth2 external redirect to %s", externalUrl)
		h.Redirect(externalUrl, http.StatusTemporaryRedirect, w, r)
	} else {
		if err = h.auth(r, provider); err != nil {
			h.renderForm(r, err)
			return
		}

		authUrl := h.getRedirectToLastURL(r)
		r.Logger().Debugf("Redirect to %s after success auth", authUrl)
		h.Redirect(authUrl, http.StatusTemporaryRedirect, w, r)
	}
}
