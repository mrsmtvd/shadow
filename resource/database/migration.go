package database

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/rubenv/sql-migrate"
	"gopkg.in/gorp.v1"
)

var idRegexp = regexp.MustCompile(`^(\d+)(.*)$`)

type ServiceMigrations interface {
	GetMigrations() migrate.MigrationSource
}

func (r *Database) GetDialect() (string, error) {
	dialect := r.GetStorage().executor.(*gorp.DbMap).Dialect

	for key := range migrate.MigrationDialects {
		if migrate.MigrationDialects[key] == dialect {
			return key, nil
		}
	}

	return "", errors.New("Unknown dialect")
}

func (r *Database) FindMigrations() ([]*migrate.Migration, error) {
	list := []*migrate.Migration{}

	for _, s := range r.application.GetServices() {
		if service, ok := s.(ServiceMigrations); ok {
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

func (r *Database) UpMigrations() (int, error) {
	dialect, err := r.GetDialect()
	if err != nil {
		return 0, err
	}

	return migrate.Exec(r.GetStorage().executor.(*gorp.DbMap).Db, dialect, r, migrate.Up)
}

func (r *Database) DownMigrations() (int, error) {
	dialect, err := r.GetDialect()
	if err != nil {
		return 0, err
	}

	return migrate.Exec(r.GetStorage().executor.(*gorp.DbMap).Db, dialect, r, migrate.Down)
}
