package database

import (
	"fmt"
	"regexp"

	"github.com/go-gorp/gorp"
	"github.com/rubenv/sql-migrate"
)

const (
	lockTimeoutInSeconds = 1
)

var idRegexp = regexp.MustCompile(`^(\d+)(.*)$`)

type hasMigrations interface {
	GetMigrations() migrate.MigrationSource
}

func (c *Component) FindMigrations() ([]*migrate.Migration, error) {
	components, err := c.application.GetComponents()
	if err != nil {
		return nil, err
	}

	list := []*migrate.Migration{}

	for _, s := range components {
		if service, ok := s.(hasMigrations); ok {
			migrations, err := service.GetMigrations().FindMigrations()

			if err != nil {
				return nil, err
			}

			for i := range migrations {
				parts := idRegexp.FindStringSubmatch(migrations[i].Id)
				if len(parts) > 2 {
					migrations[i].Id = fmt.Sprintf("%s_%s%s", parts[1], s.GetName(), parts[2])
				}
			}

			list = append(list, migrations...)
		}
	}

	return list, nil
}

func (c *Component) execWithLock(dir migrate.MigrationDirection) (int, error) {
	storage := c.GetStorage()
	dialect := storage.GetDialect()

	if dialect == "mysql" {
		lockName := c.config.GetString(ConfigMigrationsTable) + ".lock"

		result, err := storage.SelectIntByQuery("SELECT GET_LOCK(?, ?)", lockName, lockTimeoutInSeconds)
		if err != nil {
			return 0, err
		}

		if result != 1 {
			c.logger.Warn("Migrations are locked")
			return 0, nil
		}

		defer func() {
			storage.ExecByQuery("SELECT RELEASE_LOCK(?)", lockName)
		}()
	}

	return migrate.Exec(storage.executor.(*gorp.DbMap).Db, dialect, c, dir)
}

func (c *Component) UpMigrations() (int, error) {
	return c.execWithLock(migrate.Up)
}

func (c *Component) DownMigrations() (int, error) {
	return c.execWithLock(migrate.Down)
}
