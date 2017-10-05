package database

import (
	"bytes"
	"path"
	"strings"

	"github.com/rubenv/sql-migrate"
)

type AssetMigration struct {
	Source   string
	Asset    func(path string) ([]byte, error)
	AssetDir func(path string) ([]string, error)
	Dir      string
}

func (a AssetMigration) GetMigrations() []Migration {
	migrations := []Migration{}

	files, err := a.AssetDir(a.Dir)
	if err != nil {
		return nil
	}

	for _, name := range files {
		if strings.HasSuffix(name, ".sql") {
			file, err := a.Asset(path.Join(a.Dir, name))
			if err != nil {
				return nil
			}

			migration, err := migrate.ParseMigration(name, bytes.NewReader(file))
			if err != nil {
				return nil
			}

			migrations = append(migrations, NewMigration(a.Source, migration.Id, migration.Up, migration.Down, nil))
		}
	}

	return migrations
}
