package database

import (
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	GetStorage() Storage
}
