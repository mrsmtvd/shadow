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
	for _, m := range c.Migrations() {
		if m.Id() == id && m.Source() == source {
			return m
		}
	}

	return nil
}

func (c *Component) Migrations() database.Migrations {
	exists := c.getCollect()
	if len(exists) == 0 {
		return nil
	}

	migrations := make(database.Migrations, len(exists), len(exists))
	s := c.Storage().(*storage.SQL)
	executor := s.Master().(*storage.SQLExecutor)

	records, err := migrate.GetMigrationRecords(executor.DB(), s.Dialect())
	if err == nil {
		for i, m := range exists {
			var appliedAt *time.Time

			for _, record := range records {
				if record.Id == formatId(m.Source(), m.Id()) {
					appliedAt = &record.AppliedAt
					break
				}
			}

			migrations[i] = database.NewMigration(m.Source(), m.Id(), m.Up(), m.Down(), m.ModAt(), appliedAt)
		}
	} else {
		for i, m := range exists {
			migrations[i] = database.NewMigration(m.Source(), m.Id(), m.Up(), m.Down(), m.ModAt(), nil)
		}
	}

	return migrations
}

func (c *Component) getCollect() database.Migrations {
	components, err := c.application.GetComponents()
	if err != nil {
		return nil
	}

	migrations := database.Migrations{}

	for _, s := range components {
		if service, ok := s.(database.HasMigrations); ok {
			for _, migration := range service.DatabaseMigrations() {
				if !nameRegexp.MatchString(migration.Id()) {
					c.logger.Warnf("Skip migration with wrong id %s", migration.Id())
					continue
				}

				migrations = append(migrations, migration)

			}
		}
	}

	sort.Sort(migrations)

	return migrations
}

func (c *Component) FindMigrations() ([]*migrate.Migration, error) {
	migrations := []*migrate.Migration{}

	for _, m := range c.getCollect() {
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
