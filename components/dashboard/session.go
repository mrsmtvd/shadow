package dashboard

import (
	"time"

	"github.com/markbates/goth"
)

type Session interface {
	RenewToken() error
	Destroy() error
	GetString(k string) (string, error)
	PutString(k string, v string) error
	PopString(k string) (string, error)
	GetBool(k string) (bool, error)
	PutBool(k string, v bool) error
	PopBool(k string) (bool, error)
	GetInt(k string) (int, error)
	PutInt(k string, v int) error
	PopInt(k string) (int, error)
	GetInt64(k string) (int64, error)
	PutInt64(k string, v int64) error
	PopInt64(k string) (int64, error)
	GetFloat(k string) (float64, error)
	PutFloat(k string, v float64) error
	PopFloat(k string) (float64, error)
	GetTime(k string) (time.Time, error)
	PutTime(k string, v time.Time) error
	PopTime(k string) (time.Time, error)
	GetBytes(k string) ([]byte, error)
	PutBytes(k string, v []byte) error
	PopBytes(k string) ([]byte, error)
	GetObject(k string, d interface{}) error
	PutObject(k string, v interface{}) error
	PopObject(k string, d interface{}) error
	Keys() ([]string, error)
	Exists(k string) (bool, error)
	Remove(k string) error
	Clear() error
}

func SessionAuthProvider(provider goth.Provider) string {
	return ComponentName + ":" + provider.Name()
}
