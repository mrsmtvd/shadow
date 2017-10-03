package database

import (
	"github.com/kihamo/shadow"
	"github.com/rubenv/sql-migrate"
)

type Component interface {
	shadow.Component

	GetStorage() Storage

	FindMigrations() ([]*migrate.Migration, error)
}
