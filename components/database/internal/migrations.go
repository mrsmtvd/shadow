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

func (c *Component) prepareMigrations(id, source string, dir migrate.MigrationDirection) (int, error) {
	migrations := make([]*migrate.Migration, 0)
	collection := c.collection()
	limit := 0

	if id == "" && source == "" {
		for _, m := range collection {
			mig := migrate.Migration{
				Id:   formatId(m.Source(), m.Id()),
				Up:   m.Up(),
				Down: m.Down(),
			}

			migrations = append(migrations, &mig)
		}
	} else {
		s := c.Storage().(*storage.SQL)
		db := s.Master().(*storage.SQLExecutor).DB()
		dialect := s.Dialect()

		records, err := migrate.GetMigrationRecords(db, dialect)
		if err != nil {
			return -1, err
		}

		searchID := formatId(source, id)

		if len(records) == 0 {
			if dir == migrate.Up {
				for _, m := range collection {
					genID := formatId(m.Source(), m.Id())

					mig := migrate.Migration{
						Id:   genID,
						Up:   m.Up(),
						Down: m.Down(),
					}

					migrations = append(migrations, &mig)

					if genID == searchID {
						break
					}
				}
			}
		} else {
			existing := make(map[string]struct{}, len(records))
			for _, record := range records {
				existing[record.Id] = struct{}{}
			}

			if dir == migrate.Up {
				for _, m := range collection {
					genID := formatId(m.Source(), m.Id())

					mig := migrate.Migration{
						Id:   genID,
						Up:   m.Up(),
						Down: m.Down(),
					}

					migrations = append(migrations, &mig)

					if genID == searchID {
						break
					}
				}
			} else {
				for _, m := range collection {
					genID := formatId(m.Source(), m.Id())

					if genID != searchID {
						if _, ok := existing[genID]; !ok {
							continue
						}
					}

					mig := migrate.Migration{
						Id:   genID,
						Up:   m.Up(),
						Down: m.Down(),
					}

					migrations = append(migrations, &mig)

					if genID == searchID || limit > 0 {
						limit++
					}
				}
			}
		}
	}

	if len(migrations) == 0 {
		return 0, nil
	}

	return c.execWithLock(dir, migrate.MemoryMigrationSource{
		Migrations: migrations,
	}, limit)
}

func (c *Component) execWithLock(dir migrate.MigrationDirection, m migrate.MigrationSource, limit int) (int, error) {
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

	return migrate.ExecMax(executor.DB(), dialect, m, dir, limit)
}

func (c *Component) UpMigration(id, source string) error {
	n, err := c.prepareMigrations(id, source, migrate.Up)
	if err != nil {
		c.logger.Error("Up migration failed", "error", err.Error(), "source", source, "id", id)
		return err
	}

	if n == 0 {
		return errors.New("migration for source " + source + " with id " + id + " not found")
	}

	c.logger.Info("Applied migration", "source", source, "id", id)
	return nil
}

func (c *Component) UpMigrations() (n int, err error) {
	n, err = c.prepareMigrations("", "", migrate.Up)
	if err != nil {
		c.logger.Error("Up migrations failed", "error", err.Error())
	} else {
		c.logger.Infof("Applied %d migrations", n)
	}

	return n, err
}

func (c *Component) DownMigration(id, source string) error {
	n, err := c.prepareMigrations(id, source, migrate.Down)
	if err != nil {
		c.logger.Error("Down migration failed", "error", err.Error(), "source", source, "id", id)
		return err
	}

	if n == 0 {
		return errors.New("migration for source " + source + " with id " + id + " not found")
	}

	c.logger.Info("Downgraded migration", "source", source, "id", id)
	return nil
}

func (c *Component) DownMigrations() (n int, err error) {
	n, err = c.prepareMigrations("", "", migrate.Down)
	if err != nil {
		c.logger.Error("Down migrations failed", "error", err.Error())
	} else {
		c.logger.Infof("Downgraded %d migrations", n)
	}

	return n, err
}
