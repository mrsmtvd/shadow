package internal

import (
	"github.com/mrsmtvd/shadow/components/database"
)

func (c *Component) DatabaseMigrations() []database.Migration {
	fs := c.AssetFS()
	fs.Prefix = "migrations"

	return database.MigrationsFromAsset(fs)
}
