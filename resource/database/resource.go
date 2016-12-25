package database

import (
	"github.com/go-gorp/gorp"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
	"github.com/kihamo/shadow/resource/logger"
	"github.com/rubenv/sql-migrate"
)

type Resource struct {
	application *shadow.Application
	config      *config.Resource
	storage     *SqlStorage
	logger      logger.Logger
}

func (r *Resource) GetName() string {
	return "database"
}

func (r *Resource) Init(a *shadow.Application) error {
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}

	r.config = resourceConfig.(*config.Resource)

	r.application = a

	return nil
}

func (r *Resource) Run() (err error) {
	r.logger = logger.NewOrNop(r.GetName(), r.application)
	r.storage, err = NewSQLStorage(r.config.GetString(ConfigDatabaseDriver), r.config.GetString(ConfigDatabaseDsn))

	if err != nil {
		return err
	}

	r.initConns(r.config.GetInt(ConfigDatabaseMaxIdleConns), r.config.GetInt(ConfigDatabaseMaxOpenConns))
	r.initTrace(r.config.GetBool(config.ConfigDebug))

	r.storage.SetTypeConverter(TypeConverter{})

	tableName := r.config.GetString(ConfigDatabaseMigrationsTable)
	if tableName == "" {
		tableName = defaultMigrationsTableName
	}
	migrate.SetTable(tableName)

	n, err := r.UpMigrations()
	if err != nil {
		return err
	}

	r.logger.Infof("Applied %d migrations", n)

	return nil
}

func (r *Resource) GetStorage() *SqlStorage {
	return r.storage
}

func (r *Resource) initConns(idle int, open int) {
	dbMap := r.storage.executor.(*gorp.DbMap)

	dbMap.Db.SetMaxIdleConns(idle)
	dbMap.Db.SetMaxOpenConns(open)
}

func (r *Resource) initTrace(enable bool) {
	dbMap := r.storage.executor.(*gorp.DbMap)

	if enable {
		dbMap.TraceOn("", r.logger)
	} else {
		dbMap.TraceOff()
	}
}
