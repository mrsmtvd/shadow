package database

import (
	"bytes"
	"os"
	"path"
	"strings"

	"github.com/rubenv/sql-migrate"
)

func GetMigrationsFromAsset(source string, dir string, asset func(path string) ([]byte, error), assetInfo func(name string) (os.FileInfo, error), assetDir func(path string) ([]string, error)) Migrations {
	migrations := Migrations{}

	files, err := assetDir(dir)
	if err != nil {
		return nil
	}

	for _, name := range files {
		if strings.HasSuffix(name, ".sql") {
			filePath := path.Join(dir, name)

			file, err := asset(filePath)
			if err != nil {
				return nil
			}

			migration, err := migrate.ParseMigration(name, bytes.NewReader(file))
			if err != nil {
				return nil
			}

			info, err := assetInfo(filePath)
			if err != nil {
				return nil
			}

			migrations = append(migrations, NewMigration(source, migration.Id, migration.Up, migration.Down, info.ModTime(), nil))
		}
	}

	return migrations
}
