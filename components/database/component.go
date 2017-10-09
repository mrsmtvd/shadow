package database

import (
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	GetStorage() Storage
	GetMigration(id, source string) Migration
	GetAllMigrations() Migrations
}
