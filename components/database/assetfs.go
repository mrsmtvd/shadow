package database

import (
	"bytes"
	"path/filepath"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/rubenv/sql-migrate"
)

const (
	MigrationFileExt = ".sql"
)

func MigrationsFromAsset(fs *assetfs.AssetFS) []Migration {
	root, err := fs.Open("")
	if err != nil {
		return nil
	}

	files, err := root.Readdir(0)
	if err != nil {
		return nil
	}

	migrations := make([]Migration, 0)
	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != MigrationFileExt {
			continue
		}

		content, err := fs.Asset(filepath.Join(fs.Prefix, f.Name()))
		if err != nil {
			continue
		}

		migration, err := migrate.ParseMigration(f.Name(), bytes.NewReader(content))
		if err != nil {
			return nil
		}

		migrations = append(migrations, NewMigration(migration.Id, migration.Up, migration.Down, f.ModTime()))
	}

	return migrations
}
