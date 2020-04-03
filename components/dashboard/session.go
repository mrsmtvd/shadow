package dashboard

import (
	"time"

	"github.com/kihamo/shadow/components/dashboard/session"
)

type Session interface {
	FlashBag() session.FlashBag
	RenewToken() error
	Destroy() error

	GetString(key string) string
	PutString(key string, value string)
	PopString(key string) string

	GetBool(key string) bool
	PutBool(key string, value bool)
	PopBool(key string) bool

	GetInt(key string) int
	PutInt(key string, value int)
	PopInt(key string) int

	GetFloat(key string) float64
	PutFloat(key string, value float64)
	PopFloat(key string) float64

	GetTime(key string) time.Time
	PutTime(key string, value time.Time)
	PopTime(key string) time.Time

	GetBytes(key string) []byte
	PutBytes(key string, value []byte)
	PopBytes(key string) []byte

	GetObject(key string) interface{}
	PutObject(key string, value interface{})

	Keys() []string
	Exists(key string) bool
	Remove(key string)
	Clear() error
}

func AuthSessionName() string {
	return "__" + ComponentName
}
