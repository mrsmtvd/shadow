package database

import (
	"github.com/go-gorp/gorp"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
	"github.com/rubenv/sql-migrate"
)

const (
	defaultMigrationsTableName = "migrations"
)

type Database struct {
	application *shadow.Application
	config      *resource.Config
	storage     *SqlStorage
}

func (r *Database) GetName() string {
	return "database"
}

func (r *Database) GetConfigVariables() []resource.ConfigVariable {
	return []resource.ConfigVariable{
		resource.ConfigVariable{
			Key:   "database.driver",
			Value: "mysql",
			Usage: "Database driver (sqlite3, postgres, mysql, mssql and oci8)",
		},
		resource.ConfigVariable{
			Key:   "database.dsn",
			Value: "root:@tcp(localhost:3306)/shadow",
			Usage: "Database DSN",
		},
		resource.ConfigVariable{
			Key:   "database.migrations.table",
			Value: defaultMigrationsTableName,
			Usage: "Database migrations table name",
		},
		resource.ConfigVariable{
			Key:   "database.max_idle_conns",
			Value: 0,
			Usage: "Database maximum number of connections in the idle connection pool",
		},
		resource.ConfigVariable{
			Key:   "database.max_open_conns",
			Value: 0,
			Usage: "Database maximum number of connections in the idle connection pool",
		},
	}
}

func (r *Database) Init(a *shadow.Application) error {
	r.application = a
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}

	r.config = resourceConfig.(*resource.Config)

	return nil
}

func (r *Database) Run() (err error) {
	r.storage, err = NewSQLStorage(r.config.GetString("database.driver"), r.config.GetString("database.dsn"))

	if err != nil {
		return err
	}

	dbMap := r.storage.executor.(*gorp.DbMap)
	dbMap.Db.SetMaxIdleConns(r.config.GetInt("database.max_idle_conns"))
	dbMap.Db.SetMaxOpenConns(r.config.GetInt("database.max_open_conns"))

	resourceLogger, err := r.application.GetResource("logger")
	if err != nil {
		return err
	}

	logger := resourceLogger.(*resource.Logger).Get(r.GetName())

	if r.config.GetBool("debug") {
		dbMap.TraceOn("", logger)
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

	logger.Debugf("Applied %d migrations", n)

	return nil
}

func (r *Database) GetStorage() *SqlStorage {
	return r.storage
}
