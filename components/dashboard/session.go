package dashboard

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/engine/memstore"
	"github.com/alexedwards/scs/session"
)

const (
	SessionUsername = "username"
)

type Session struct {
	response http.ResponseWriter
	request  *http.Request
}

func NewSession(w http.ResponseWriter, r *http.Request) *Session {
	return &Session{
		response: w,
		request:  r,
	}
}

func NewSessionManager() func(h http.Handler) http.Handler {
	session.CookieName = "shadow.token"

	sessionEngine := memstore.New(0)
	return session.Manage(sessionEngine)
}

func (s *Session) RegenerateToken() error {
	return session.RegenerateToken(s.request)
}

func (s *Session) Renew() error {
	return session.Renew(s.request)
}

func (s *Session) Destroy() error {
	return session.Destroy(s.response, s.request)
}

func (s *Session) GetString(k string) (string, error) {
	return session.GetString(s.request, k)
}

func (s *Session) PutString(k string, v string) error {
	return session.PutString(s.request, k, v)
}

func (s *Session) PopString(k string) (string, error) {
	return session.PopString(s.request, k)
}

func (s *Session) GetBool(k string) (bool, error) {
	return session.GetBool(s.request, k)
}

func (s *Session) PutBool(k string, v bool) error {
	return session.PutBool(s.request, k, v)
}

func (s *Session) PopBool(k string) (bool, error) {
	return session.PopBool(s.request, k)
}

func (s *Session) GetInt(k string) (int, error) {
	return session.GetInt(s.request, k)
}

func (s *Session) PutInt(k string, v int) error {
	return session.PutInt(s.request, k, v)
}

func (s *Session) PopInt(k string) (int, error) {
	return session.PopInt(s.request, k)
}

func (s *Session) GetInt64(k string) (int64, error) {
	return session.GetInt64(s.request, k)
}

func (s *Session) PutInt64(k string, v int64) error {
	return session.PutInt64(s.request, k, v)
}

func (s *Session) PopInt64(k string) (int64, error) {
	return session.PopInt64(s.request, k)
}

func (s *Session) GetFloat(k string) (float64, error) {
	return session.GetFloat(s.request, k)
}

func (s *Session) PutFloat(k string, v float64) error {
	return session.PutFloat(s.request, k, v)
}

func (s *Session) PopFloat(k string) (float64, error) {
	return session.PopFloat(s.request, k)
}

func (s *Session) GetTime(k string) (time.Time, error) {
	return session.GetTime(s.request, k)
}

func (s *Session) PutTime(k string, v time.Time) error {
	return session.PutTime(s.request, k, v)
}

func (s *Session) PopTime(k string) (time.Time, error) {
	return session.PopTime(s.request, k)
}

func (s *Session) GetBytes(k string) ([]byte, error) {
	return session.GetBytes(s.request, k)
}

func (s *Session) PutBytes(k string, v []byte) error {
	return session.PutBytes(s.request, k, v)
}

func (s *Session) PopBytes(k string) ([]byte, error) {
	return session.PopBytes(s.request, k)
}

func (s *Session) GetObject(k string, d interface{}) error {
	return session.GetObject(s.request, k, d)
}

func (s *Session) PutObject(k string, v interface{}) error {
	return session.PutObject(s.request, k, v)
}

func (s *Session) PopObject(k string, d interface{}) error {
	return session.PopObject(s.request, k, d)
}

func (s *Session) Keys() ([]string, error) {
	return session.Keys(s.request)
}

func (s *Session) Exists(k string) (bool, error) {
	return session.Exists(s.request, k)
}

func (s *Session) Remove(k string) error {
	return session.Remove(s.request, k)
}

func (s *Session) Clear() error {
	return session.Clear(s.request)
}
