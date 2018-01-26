package internal

import (
	"strings"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/database"
	"github.com/kihamo/shadow/components/database/balancer"
	"github.com/kihamo/shadow/components/database/storage"
	"github.com/kihamo/shadow/components/logger"
	"github.com/rubenv/sql-migrate"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logger.Logger
	routes      []dashboard.Route
	storage     database.Storage
}

func (c *Component) GetName() string {
	return database.ComponentName
}

func (c *Component) GetVersion() string {
	return database.ComponentVersion
}

func (c *Component) GetDependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: logger.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	return nil
}

func (c *Component) Run() (err error) {
	c.logger = logger.NewOrNop(c.GetName(), c.application)

	var slaves []string
	if slavesFromConfig := c.config.GetString(database.ConfigDsnSlaves); slavesFromConfig != "" {
		slaves = strings.Split(slavesFromConfig, ";")
	}

	s, err := storage.NewSQL(
		c.config.GetString(database.ConfigDriver),
		c.config.GetString(database.ConfigDsnMaster),
		slaves,
		map[string]string{
			storage.DialectOptionEngine:   c.config.GetString(database.ConfigDialectEngine),
			storage.DialectOptionEncoding: c.config.GetString(database.ConfigDialectEncoding),
			storage.DialectOptionVersion:  c.config.GetString(database.ConfigDialectVersion),
		},
		c.config.GetBool(database.ConfigAllowUseMasterAsSlave),
	)

	if err != nil {
		return err
	}

	s.SetMaxOpenConns(c.config.GetInt(database.ConfigMaxOpenConns))
	s.SetMaxIdleConns(c.config.GetInt(database.ConfigMaxIdleConns))
	s.SetConnMaxLifetime(c.config.GetDuration(database.ConfigConnMaxLifetime))

	s.SetTypeConverter(TypeConverter{})
	c.storage = s

	c.initTrace(s, c.config.GetBool(config.ConfigDebug))
	c.initBalancer(s, c.config.GetString(database.ConfigBalancer))

	migrate.SetSchema(c.config.GetString(database.ConfigMigrationsSchema))
	migrate.SetTable(c.config.GetString(database.ConfigMigrationsTable))

	n, err := c.UpMigrations()
	if err != nil {
		return err
	}

	c.logger.Infof("Applied %d migrations", n)

	return nil
}

func (c *Component) initTrace(s database.Storage, d bool) {
	if d {
		s.(*storage.SQL).TraceOn(c.logger)
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
	return c.storage
}
