package internal

import (
	"strings"
	"sync"

	"github.com/mrsmtvd/shadow"
	"github.com/mrsmtvd/shadow/components/config"
	"github.com/mrsmtvd/shadow/components/database"
	"github.com/mrsmtvd/shadow/components/database/balancer"
	"github.com/mrsmtvd/shadow/components/database/storage"
	"github.com/mrsmtvd/shadow/components/i18n"
	"github.com/mrsmtvd/shadow/components/logging"
	"github.com/mrsmtvd/shadow/components/metrics"
	migrate "github.com/rubenv/sql-migrate"
)

type Component struct {
	mutex sync.RWMutex

	application shadow.Application
	config      config.Component
	logger      logging.Logger
	storage     database.Storage

	migrationsIsUp  bool
	migrationsError error
}

func (c *Component) Name() string {
	return database.ComponentName
}

func (c *Component) Version() string {
	return database.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: i18n.ComponentName,
		},
		{
			Name: logging.ComponentName,
		},
		{
			Name: metrics.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) error {
	c.logger = logging.DefaultLazyLogger(c.Name())

	<-a.ReadyComponent(config.ComponentName)

	var slaves []string
	if slavesFromConfig := c.config.String(database.ConfigDsnSlaves); slavesFromConfig != "" {
		slaves = strings.Split(slavesFromConfig, ";")

		for i := 0; i < len(slaves); i++ {
			slaves[i] = strings.TrimSpace(slaves[i])

			if slaves[i] == "" {
				slaves = append(slaves[:i], slaves[i+1:]...)
			}
		}
	}

	s, err := storage.NewSQL(
		c.config.String(database.ConfigDriver),
		c.config.String(database.ConfigDsnMaster),
		slaves,
		map[string]string{
			storage.DialectOptionEngine:   c.config.String(database.ConfigDialectEngine),
			storage.DialectOptionEncoding: c.config.String(database.ConfigDialectEncoding),
			storage.DialectOptionVersion:  c.config.String(database.ConfigDialectVersion),
		},
		c.config.Bool(database.ConfigAllowUseMasterAsSlave),
	)

	if err != nil {
		return err
	}

	s.SetMaxOpenConns(c.config.Int(database.ConfigMaxOpenConns))
	s.SetMaxIdleConns(c.config.Int(database.ConfigMaxIdleConns))
	s.SetConnMaxLifetime(c.config.Duration(database.ConfigConnMaxLifetime))

	s.SetTypeConverter(TypeConverter{})

	c.mutex.Lock()
	c.storage = s
	c.mutex.Unlock()

	c.initTrace(s, c.config.Bool(config.ConfigDebug))
	c.initBalancer(s, c.config.String(database.ConfigBalancer))

	migrate.SetSchema(c.config.String(database.ConfigMigrationsSchema))
	migrate.SetTable(c.config.String(database.ConfigMigrationsTable))

	ready <- struct{}{}

	_, err = c.UpMigrations()

	c.mutex.Lock()
	c.migrationsIsUp = err == nil
	c.migrationsError = err
	c.mutex.Unlock()

	return nil
}

func (c *Component) initTrace(s database.Storage, d bool) {
	if d {
		s.(*storage.SQL).TraceOn(&logger{c.logger})
	} else {
		s.(*storage.SQL).TraceOff()
	}
}

func (c *Component) initBalancer(s database.Storage, b string) {
	switch b {
	case database.BalancerRandom:
		s.SetBalancer(balancer.NewRandom())
	default:
		s.SetBalancer(balancer.NewRoundRobin())
	}
}

func (c *Component) Storage() database.Storage {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.storage
}
