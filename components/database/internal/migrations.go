package internal

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/kihamo/shadow/components/database"
	"github.com/kihamo/shadow/components/database/storage"
	"github.com/rubenv/sql-migrate"
)

const (
	lockTimeoutInSeconds = 1
)

var nameRegexp = regexp.MustCompile(`^([1-9]\d{3}[0-1]\d[0-3]\d[0-2]\d[0-5]\d[0-5]\d)(.*)$`)

func formatId(source, id string) string {
	parts := nameRegexp.FindStringSubmatch(id)
	if len(parts) > 2 {
		id = fmt.Sprintf("%s_%s%s", parts[1], source, parts[2])
	}

	return id
}

func (c *Component) Migration(id, source string) database.Migration {
	for _, m := range c.collection() {
		if m.Id() == id && m.Source() == source {
			return m
		}
	}

	return nil
}

func (c *Component) Migrations() []database.Migration {
	collection := c.collection()
	if len(collection) == 0 {
		return nil
	}

	migrations := make([]database.Migration, len(collection), len(collection))
	s := c.Storage().(*storage.SQL)
	executor := s.Master().(*storage.SQLExecutor)

	records, err := migrate.GetMigrationRecords(executor.DB(), s.Dialect())
	if err == nil {
		for i, m := range collection {
			var appliedAt *time.Time

			for _, record := range records {
				if record.Id == formatId(m.Source(), m.Id()) {
					appliedAt = &record.AppliedAt
					break
				}
			}

			m.SetAppliedAt(appliedAt)
			migrations[i] = m
		}
	} else {
		for i, m := range collection {
			migrations[i] = m
		}
	}

	return migrations
}

func (c *Component) collection() MigrationsCollection {
	components, err := c.application.GetComponents()
	if err != nil {
		return nil
	}

	migrations := MigrationsCollection{}

	for _, component := range components {
		if componentMigrations, ok := component.(database.HasMigrations); ok {
			for _, migration := range componentMigrations.DatabaseMigrations() {
				if !nameRegexp.MatchString(migration.Id()) {
					c.logger.Warnf("Skip migration with wrong id %s", migration.Id())
					continue
				}

				migrations = append(migrations, NewMigrationItem(migration, component.Name()))

			}
		}
	}

	sort.Sort(migrations)

	return migrations
}

func (c *Component) FindMigrations() ([]*migrate.Migration, error) {
	collection := c.collection()
	migrations := make([]*migrate.Migration, 0, len(collection))

	for _, m := range collection {
		mig := migrate.Migration{
			Id:   formatId(m.Source(), m.Id()),
			Up:   m.Up(),
			Down: m.Down(),
		}

		migrations = append(migrations, &mig)
		c.logger.Debugf("Found %s migration and converted to %s", m.Id(), mig.Id)
	}

	return migrations, nil
}

func (c *Component) execWithLock(dir migrate.MigrationDirection) (int, error) {
	if c.Storage() == nil {
		return -1, errors.New("Storage isn't initialized")
	}

	s := c.Storage().(*storage.SQL)
	executor := s.Master().(*storage.SQLExecutor)
	dialect := s.Dialect()

	if dialect == storage.DialectMySQL {
		lockName := c.config.String(database.ConfigMigrationsTable) + ".lock"

		result, err := executor.SelectIntByQuery("SELECT GET_LOCK(?, ?)", lockName, lockTimeoutInSeconds)
		if err != nil {
			return 0, err
		}

		if result != 1 {
			c.logger.Warn("Migrations are locked")
			return 0, nil
		}

		defer func() {
			executor.ExecByQuery("SELECT RELEASE_LOCK(?)", lockName)
		}()
	}

	return migrate.Exec(executor.DB(), dialect, c, dir)
}

func (c *Component) UpMigrations() (int, error) {
	return c.execWithLock(migrate.Up)
}

func (c *Component) DownMigrations() (int, error) {
	return c.execWithLock(migrate.Down)
}
