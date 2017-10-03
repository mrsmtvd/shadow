package internal

import (
	"sync"

	"github.com/go-gorp/gorp"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/database"
	"github.com/kihamo/shadow/components/logger"
	"github.com/rubenv/sql-migrate"
)

type Component struct {
	application shadow.Application
	config      config.Component
	storage     *SqlStorage
	logger      logger.Logger
	routes      []dashboard.Route

	mutex sync.Mutex
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
	c.storage, err = NewSQLStorage(
		c.config.GetString(database.ConfigDriver),
		c.config.GetString(database.ConfigDsn),
		map[string]string{
			DialectOptionEngine:   c.config.GetString(database.ConfigDialectEngine),
			DialectOptionEncoding: c.config.GetString(database.ConfigDialectEncoding),
			DialectOptionVersion:  c.config.GetString(database.ConfigDialectVersion),
		},
	)

	if err != nil {
		return err
	}

	c.initConns(c.config.GetInt(database.ConfigMaxIdleConns), c.config.GetInt(database.ConfigMaxOpenConns))
	c.initTrace(c.config.GetBool(config.ConfigDebug))

	c.storage.SetTypeConverter(TypeConverter{})

	tableName := c.config.GetString(database.ConfigMigrationsTable)
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

func (c *Component) GetStorage() database.Storage {
	return c.storage
}

func (c *Component) initConns(idle int, open int) {
	dbMap := c.storage.executor.(*gorp.DbMap)

	dbMap.Db.SetMaxIdleConns(idle)
	dbMap.Db.SetMaxOpenConns(open)
}

func (c *Component) initTrace(enable bool) {
	dbMap := c.storage.executor.(*gorp.DbMap)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if enable {
		dbMap.TraceOn("", c.logger)
	} else {
		dbMap.TraceOff()
	}
}
