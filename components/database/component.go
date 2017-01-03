package database

import (
	"github.com/go-gorp/gorp"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
	"github.com/rubenv/sql-migrate"
)

type Component struct {
	application *shadow.Application
	config      *config.Component
	storage     *SqlStorage
	logger      logger.Logger
}

func (c *Component) GetName() string {
	return "database"
}

func (c *Component) GetVersion() string {
	return "1.0.0"
}

func (c *Component) Init(a *shadow.Application) error {
	resourceConfig, err := a.GetComponent("config")
	if err != nil {
		return err
	}

	c.config = resourceConfig.(*config.Component)

	c.application = a

	return nil
}

func (c *Component) Run() (err error) {
	c.logger = logger.NewOrNop(c.GetName(), c.application)
	c.storage, err = NewSQLStorage(c.config.GetString(ConfigDatabaseDriver), c.config.GetString(ConfigDatabaseDsn))

	if err != nil {
		return err
	}

	c.initConns(c.config.GetInt(ConfigDatabaseMaxIdleConns), c.config.GetInt(ConfigDatabaseMaxOpenConns))
	c.initTrace(c.config.GetBool(config.ConfigDebug))

	c.storage.SetTypeConverter(TypeConverter{})

	tableName := c.config.GetString(ConfigDatabaseMigrationsTable)
	if tableName == "" {
		tableName = defaultMigrationsTableName
	}
	migrate.SetTable(tableName)

	n, err := c.UpMigrations()
	if err != nil {
		return err
	}

	c.logger.Infof("Applied %d migrations", n)

	return nil
}

func (c *Component) GetStorage() *SqlStorage {
	return c.storage
}

func (c *Component) initConns(idle int, open int) {
	dbMap := c.storage.executor.(*gorp.DbMap)

	dbMap.Db.SetMaxIdleConns(idle)
	dbMap.Db.SetMaxOpenConns(open)
}

func (c *Component) initTrace(enable bool) {
	dbMap := c.storage.executor.(*gorp.DbMap)

	if enable {
		dbMap.TraceOn("", c.logger)
	} else {
		dbMap.TraceOff()
	}
}
