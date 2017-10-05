package internal

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/kihamo/shadow/components/database"
	"github.com/rubenv/sql-migrate"
)

const (
	lockTimeoutInSeconds = 1
)

var idRegexp = regexp.MustCompile(`^(\d+)(.*)$`)

type migrations []database.Migration

func (m migrations) Len() int {
	return len(m)
}

func (m migrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m migrations) Less(i, j int) bool {
	bySource := strings.Compare(m[i].Source(), m[j].Source())

	if bySource == 0 {
		return strings.Compare(m[i].Id(), m[j].Id()) < 0
	}

	return bySource < 0
}

func formatId(source, id string) string {
	parts := idRegexp.FindStringSubmatch(id)
	if len(parts) > 2 {
		id = fmt.Sprintf("%s_%s%s", parts[1], source, parts[2])
	}

	return id
}

func (c *Component) GetMigration(id, source string) database.Migration {
	for _, m := range c.GetAllMigrations() {
		if m.Id() == id && m.Source() == source {
			return m
		}
	}

	return nil
}

func (c *Component) GetAllMigrations() []database.Migration {
	exists := c.getCollect()
	if len(exists) == 0 {
		return nil
	}

	migrations := make([]database.Migration, len(exists), len(exists))
	storage := c.GetStorage()

	records, err := migrate.GetMigrationRecords(storage.(*SqlStorage).executor.(*gorp.DbMap).Db, storage.GetDialect())
	if err == nil {
		for i, m := range exists {
			var appliedAt *time.Time

			for _, record := range records {
				if record.Id == formatId(m.Source(), m.Id()) {
					appliedAt = &record.AppliedAt
					break
				}
			}

			migrations[i] = database.NewMigration(m.Source(), m.Id(), m.Up(), m.Down(), appliedAt)
		}
	} else {
		for i, m := range exists {
			migrations[i] = database.NewMigration(m.Source(), m.Id(), m.Up(), m.Down(), nil)
		}
	}

	return migrations
}

func (c *Component) getCollect() []database.Migration {
	components, err := c.application.GetComponents()
	if err != nil {
		return nil
	}

	migrations := migrations{}

	for _, s := range components {
		if service, ok := s.(database.HasMigrations); ok {
			for _, migration := range service.GetMigrations() {
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
	}

	return migrations, nil
}

func (c *Component) execWithLock(dir migrate.MigrationDirection) (int, error) {
	storage := c.GetStorage()
	dialect := storage.GetDialect()

	if dialect == DialectMySQL {
		lockName := c.config.GetString(database.ConfigMigrationsTable) + ".lock"

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

	return migrate.Exec(storage.(*SqlStorage).executor.(*gorp.DbMap).Db, dialect, c, dir)
}

func (c *Component) UpMigrations() (int, error) {
	return c.execWithLock(migrate.Up)
}

func (c *Component) DownMigrations() (int, error) {
	return c.execWithLock(migrate.Down)
}
