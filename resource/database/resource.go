package database

import (
	"github.com/go-gorp/gorp"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource/config"
	"github.com/kihamo/shadow/resource/logger"
	"github.com/rs/xlog"
	"github.com/rubenv/sql-migrate"
)

const (
	defaultMigrationsTableName = "migrations"
)

type Resource struct {
	application *shadow.Application
	config      *config.Resource
	storage     *SqlStorage
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
	r.storage, err = NewSQLStorage(r.config.GetString("database.driver"), r.config.GetString("database.dsn"))

	if err != nil {
		return err
	}

	dbMap := r.storage.executor.(*gorp.DbMap)
	dbMap.Db.SetMaxIdleConns(r.config.GetInt("database.max_idle_conns"))
	dbMap.Db.SetMaxOpenConns(r.config.GetInt("database.max_open_conns"))

	var l xlog.Logger
	if resourceLogger, err := r.application.GetResource("logger"); err == nil {
		l = resourceLogger.(*logger.Resource).Get(r.GetName())
	} else {
		l = xlog.NopLogger
	}

	if r.config.GetBool("debug") {
		dbMap.TraceOn("", newDatabaseLogger(l))
	}

	r.storage.SetTypeConverter(TypeConverter{})

	tableName := r.config.GetString("database.migrations.table")
	if tableName == "" {
		tableName = defaultMigrationsTableName
	}
	migrate.SetTable(tableName)

	n, err := r.UpMigrations()
	if err != nil {
		return err
	}

	l.Infof("Applied %d migrations", n)

	return nil
}

func (r *Resource) GetStorage() *SqlStorage {
	return r.storage
}
