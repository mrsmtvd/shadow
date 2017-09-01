package dashboard

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs"
)

const (
	SessionUsername = "username"
)

type Session struct {
	session  *scs.Session
	response http.ResponseWriter
}

func NewSession(s *scs.Session, w http.ResponseWriter) *Session {
	return &Session{
		session:  s,
		response: w,
	}
}

func (s *Session) RenewToken() error {
	return s.session.RenewToken(s.response)
}

func (s *Session) Destroy() error {
	return s.session.Destroy(s.response)
}

func (s *Session) GetString(k string) (string, error) {
	return s.session.GetString(k)
}

func (s *Session) PutString(k string, v string) error {
	return s.session.PutString(s.response, k, v)
}

func (s *Session) PopString(k string) (string, error) {
	return s.session.PopString(s.response, k)
}

func (s *Session) GetBool(k string) (bool, error) {
	return s.session.GetBool(k)
}

func (s *Session) PutBool(k string, v bool) error {
	return s.session.PutBool(s.response, k, v)
}

func (s *Session) PopBool(k string) (bool, error) {
	return s.session.PopBool(s.response, k)
}

func (s *Session) GetInt(k string) (int, error) {
	return s.session.GetInt(k)
}

func (s *Session) PutInt(k string, v int) error {
	return s.session.PutInt(s.response, k, v)
}

func (s *Session) PopInt(k string) (int, error) {
	return s.session.PopInt(s.response, k)
}

func (s *Session) GetInt64(k string) (int64, error) {
	return s.session.GetInt64(k)
}

func (s *Session) PutInt64(k string, v int64) error {
	return s.session.PutInt64(s.response, k, v)
}

func (s *Session) PopInt64(k string) (int64, error) {
	return s.session.PopInt64(s.response, k)
}

func (s *Session) GetFloat(k string) (float64, error) {
	return s.session.GetFloat(k)
}

func (s *Session) PutFloat(k string, v float64) error {
	return s.session.PutFloat(s.response, k, v)
}

func (s *Session) PopFloat(k string) (float64, error) {
	return s.session.PopFloat(s.response, k)
}

func (s *Session) GetTime(k string) (time.Time, error) {
	return s.session.GetTime(k)
}

func (s *Session) PutTime(k string, v time.Time) error {
	return s.session.PutTime(s.response, k, v)
}

func (s *Session) PopTime(k string) (time.Time, error) {
	return s.session.PopTime(s.response, k)
}

func (s *Session) GetBytes(k string) ([]byte, error) {
	return s.session.GetBytes(k)
}

func (s *Session) PutBytes(k string, v []byte) error {
	return s.session.PutBytes(s.response, k, v)
}

func (s *Session) PopBytes(k string) ([]byte, error) {
	return s.session.PopBytes(s.response, k)
}

func (s *Session) GetObject(k string, d interface{}) error {
	return s.session.GetObject(k, d)
}

func (s *Session) PutObject(k string, v interface{}) error {
	return s.session.PutObject(s.response, k, v)
}

func (s *Session) PopObject(k string, d interface{}) error {
	return s.session.PopObject(s.response, k, d)
}

func (s *Session) Keys() ([]string, error) {
	return s.session.Keys()
}

func (s *Session) Exists(k string) (bool, error) {
	return s.session.Exists(k)
}

func (s *Session) Remove(k string) error {
	return s.session.Remove(s.response, k)
}

func (s *Session) Clear() error {
	return s.session.Clear(s.response)
}
