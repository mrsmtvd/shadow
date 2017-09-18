package password

import (
	"errors"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

type Provider struct {
	providerName string
	authURL      string
	users        map[string]string
}

func New(authURL string, users map[string]string) *Provider {
	return &Provider{
		providerName: "password",
		authURL:      authURL,
		users:        users,
	}
}

func (p *Provider) Name() string {
	return p.providerName
}

func (p *Provider) SetName(name string) {
	p.providerName = name
}

func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	return &Session{
		AuthURL: p.authURL,
		State:   state,
	}, nil
}

func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	sess := session.(*Session)
	user := goth.User{
		UserID:   sess.Username,
		Name:     sess.Username,
		Provider: p.Name(),
	}

	if !sess.Valid {
		return user, errors.New("User isn't authorized")
	}

	return user, nil
}

func (p *Provider) Debug(debug bool) {}

func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	return nil, nil
}

func (p *Provider) RefreshTokenAvailable() bool {
	return false
}
