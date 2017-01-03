package database

import (
	"fmt"
	"regexp"

	"github.com/go-gorp/gorp"
	"github.com/rubenv/sql-migrate"
)

var idRegexp = regexp.MustCompile(`^(\d+)(.*)$`)

type hasMigrations interface {
	GetMigrations() migrate.MigrationSource
}

func (c *Component) FindMigrations() ([]*migrate.Migration, error) {
	list := []*migrate.Migration{}

	for _, s := range c.application.GetComponents() {
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

func (c *Component) UpMigrations() (int, error) {
	dialect, err := c.GetStorage().GetDialect()
	if err != nil {
		return 0, err
	}

	return migrate.Exec(c.GetStorage().executor.(*gorp.DbMap).Db, dialect, c, migrate.Up)
}

func (c *Component) DownMigrations() (int, error) {
	dialect, err := c.GetStorage().GetDialect()
	if err != nil {
		return 0, err
	}

	return migrate.Exec(c.GetStorage().executor.(*gorp.DbMap).Db, dialect, c, migrate.Down)
}
