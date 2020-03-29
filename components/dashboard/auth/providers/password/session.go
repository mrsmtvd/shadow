package password

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/markbates/goth"
)

// Session stores data during the auth process with Gitlab.
type Session struct {
	AuthURL  string
	State    string
	Username string
	Password string
	Valid    bool
}

var _ goth.Session = &Session{}

func (s Session) GetAuthURL() (string, error) {
	u, err := url.Parse(s.AuthURL)
	if err != nil {
		return "", err
	}

	values := u.Query()
	values.Add("state", s.State)
	u.RawQuery = values.Encode()

	return u.String(), nil
}

func (s *Session) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	p := provider.(*Provider)
	username := params.Get("username")

	userPassword, ok := p.users[username]
	if !ok {
		return "", errors.New("Invalid username")
	}

	password := params.Get("password")

	if userPassword != password {
		return "", errors.New("Invalid password")
	}

	s.Username = username
	s.Password = password
	s.Valid = true

	return "", nil
}

func (s Session) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s Session) String() string {
	return s.Marshal()
}

func (p *Provider) UnmarshalSession(data string) (goth.Session, error) {
	s := &Session{}
	err := json.NewDecoder(strings.NewReader(data)).Decode(s)

	return s, err
}
