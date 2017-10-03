package database

import (
	"github.com/rubenv/sql-migrate"
)

type HasMigrations interface {
	GetMigrations() migrate.MigrationSource
}
