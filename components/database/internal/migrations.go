package internal

import (
	"errors"
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
		id = parts[1] + "_" + source + parts[2]
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

	migrations := make([]database.Migration, len(collection))
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
					c.logger.Warn("Skip migration with wrong id " + migration.Id())
					continue
				}

				migrations = append(migrations, NewMigrationItem(migration, component.Name()))
			}
		}
	}

	sort.Sort(migrations)

	return migrations
}

func (c *Component) prepareMigration(id, source string) *migrate.Migration {
	for _, m := range c.collection() {
		if m.Id() == id && m.Source() == source {
			return &migrate.Migration{
				Id:   formatId(m.Source(), m.Id()),
				Up:   m.Up(),
				Down: m.Down(),
			}
		}
	}

	return nil
}

func (c *Component) prepareMigrations() []*migrate.Migration {
	collection := c.collection()
	migrations := make([]*migrate.Migration, 0, len(collection))

	for _, m := range collection {
		mig := migrate.Migration{
			Id:   formatId(m.Source(), m.Id()),
			Up:   m.Up(),
			Down: m.Down(),
		}

		migrations = append(migrations, &mig)
		c.logger.Debug("Found " + m.Id() + " migration and converted to " + mig.Id)
	}

	return migrations
}

func (c *Component) execWithLock(dir migrate.MigrationDirection, m migrate.MigrationSource) (int, error) {
	if len(c.Migrations()) == 0 {
		return 0, nil
	}

	if c.Storage() == nil {
		return -1, errors.New("storage isn't initialized")
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
			_, _ = executor.ExecByQuery("SELECT RELEASE_LOCK(?)", lockName)
		}()
	}

	return migrate.Exec(executor.DB(), dialect, m, dir)
}

func (c *Component) UpMigration(id, source string) error {
	m := c.prepareMigration(id, source)
	if m == nil {
		return errors.New("migration for source " + source + " with id " + id + " not found")
	}

	_, err := c.execWithLock(migrate.Up, migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{m},
	})

	if err != nil {
		c.logger.Error("Up migration failed", "error", err.Error(), "source", source, "id", id)
	} else {
		c.logger.Info("Applied migration", "source", source, "id", id)
	}

	return err
}

func (c *Component) UpMigrations() (n int, err error) {
	n, err = c.execWithLock(migrate.Up, migrate.MemoryMigrationSource{
		Migrations: c.prepareMigrations(),
	})
	if err != nil {
		c.logger.Error("Up migrations failed", "error", err.Error())
	} else {
		c.logger.Infof("Applied %d migrations", n)
	}

	return
}

func (c *Component) DownMigration(id, source string) error {
	m := c.prepareMigration(id, source)
	if m == nil {
		return errors.New("migration for source " + source + " with id " + id + " not found")
	}

	_, err := c.execWithLock(migrate.Down, migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{m},
	})

	if err != nil {
		c.logger.Error("Down migration failed", "error", err.Error(), "source", source, "id", id)
	} else {
		c.logger.Info("Downgraded migration", "source", source, "id", id)
	}

	return err
}

func (c *Component) DownMigrations() (n int, err error) {
	n, err = c.execWithLock(migrate.Down, migrate.MemoryMigrationSource{
		Migrations: c.prepareMigrations(),
	})

	if err != nil {
		c.logger.Error("Down migrations failed", "error", err.Error())
	} else {
		c.logger.Infof("Downgraded %d migrations", n)
	}

	return
}
