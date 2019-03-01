package internal

import (
	"github.com/kihamo/shadow/components/database"
)

func (c *Component) DatabaseMigrations() []database.Migration {
	fs := c.AssetFS()
	fs.Prefix = "migrations"

	return database.MigrationsFromAsset(fs)
}
