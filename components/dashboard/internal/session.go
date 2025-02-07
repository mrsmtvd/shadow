package internal

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mrsmtvd/shadow/components/dashboard/session"
)

const (
	SessionFlashBag = "flash-bag"
)

type Session struct {
	session  *scs.SessionManager
	flashBag *session.AutoExpireFlashBag
	request  *http.Request
}

func NewSession(s *scs.SessionManager, request *http.Request) *Session {
	r := &Session{
		session:  s,
		flashBag: session.NewAutoExpireFlashBag(),
		request:  request,
	}
	r.init()

	return r
}

func (s *Session) init() {
	if !s.Exists(SessionFlashBag) {
		content, _ := json.Marshal(nil)
		s.PutBytes(SessionFlashBag, content)

		return
	}

	dump := make(map[string][]string)

	if err := json.Unmarshal(s.GetBytes(SessionFlashBag), &dump); err != nil {
		return
	}

	for level, messages := range dump {
		for _, message := range messages {
			s.flashBag.Add(level, message)
		}
	}

	s.flashBag.Commit()
}

func (s *Session) Flush() {
	if s.flashBag.Changed() {
		if content, err := json.Marshal(s.flashBag.All()); err == nil {
			s.PutBytes(SessionFlashBag, content)
		}
	}
}

func (s *Session) FlashBag() session.FlashBag {
	return s.flashBag
}

func (s *Session) RenewToken() error {
	return s.session.RenewToken(s.request.Context())
}

func (s *Session) Destroy() error {
	return s.session.Destroy(s.request.Context())
}

func (s *Session) GetString(key string) string {
	return s.session.GetString(s.request.Context(), key)
}

func (s *Session) PutString(key string, value string) {
	s.session.Put(s.request.Context(), key, value)
}

func (s *Session) PopString(key string) string {
	return s.session.PopString(s.request.Context(), key)
}

func (s *Session) GetBool(key string) bool {
	return s.session.GetBool(s.request.Context(), key)
}

func (s *Session) PutBool(key string, value bool) {
	s.session.Put(s.request.Context(), key, value)
}

func (s *Session) PopBool(key string) bool {
	return s.session.PopBool(s.request.Context(), key)
}

func (s *Session) GetInt(key string) int {
	return s.session.GetInt(s.request.Context(), key)
}

func (s *Session) PutInt(key string, value int) {
	s.session.Put(s.request.Context(), key, value)
}

func (s *Session) PopInt(key string) int {
	return s.session.PopInt(s.request.Context(), key)
}

func (s *Session) GetFloat(key string) float64 {
	return s.session.GetFloat(s.request.Context(), key)
}

func (s *Session) PutFloat(key string, value float64) {
	s.session.Put(s.request.Context(), key, value)
}

func (s *Session) PopFloat(key string) float64 {
	return s.session.PopFloat(s.request.Context(), key)
}

func (s *Session) GetTime(key string) time.Time {
	return s.session.GetTime(s.request.Context(), key)
}

func (s *Session) PutTime(key string, value time.Time) {
	s.session.Put(s.request.Context(), key, value)
}

func (s *Session) PopTime(key string) time.Time {
	return s.session.PopTime(s.request.Context(), key)
}

func (s *Session) GetBytes(key string) []byte {
	return s.session.GetBytes(s.request.Context(), key)
}

func (s *Session) PutBytes(key string, value []byte) {
	s.session.Put(s.request.Context(), key, value)
}

func (s *Session) PopBytes(key string) []byte {
	return s.session.PopBytes(s.request.Context(), key)
}

func (s *Session) GetObject(key string) interface{} {
	return s.session.Get(s.request.Context(), key)
}

func (s *Session) PutObject(key string, value interface{}) {
	s.session.Put(s.request.Context(), key, value)
}

func (s *Session) Keys() []string {
	return s.session.Keys(s.request.Context())
}

func (s *Session) Exists(key string) bool {
	return s.session.Exists(s.request.Context(), key)
}

func (s *Session) Remove(key string) {
	s.session.Remove(s.request.Context(), key)
}

func (s *Session) Clear() error {
	return s.session.Clear(s.request.Context())
}
