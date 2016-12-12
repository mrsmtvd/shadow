package database

import (
	"github.com/go-gorp/gorp"
	"github.com/kihamo/shadow/resource/config"
	"github.com/rubenv/sql-migrate"
)

const (
	ConfigDatabaseDriver          = "database.driver"
	ConfigDatabaseDsn             = "database.dsn"
	ConfigDatabaseMigrationsTable = "database.migrations.table"
	ConfigDatabaseMaxIdleConns    = "database.max_idle_conns"
	ConfigDatabaseMaxOpenConns    = "database.max_open_conns"

	defaultMigrationsTableName = "migrations"
)

func (r *Resource) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:     ConfigDatabaseDriver,
			Default: "mysql",
			Usage:   "Database driver (sqlite3, postgres, mysql, mssql and oci8)",
			Type:    config.ValueTypeString,
		},
		{
			Key:   ConfigDatabaseDsn,
			Usage: "Database DSN",
			Type:  config.ValueTypeString,
		},
		{
			Key:      ConfigDatabaseMigrationsTable,
			Default:  defaultMigrationsTableName,
			Usage:    "Database migrations table name",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigDatabaseMaxIdleConns,
			Default:  0,
			Usage:    "Database maximum number of connections in the idle connection pool",
			Type:     config.ValueTypeInt,
			Editable: true,
		},
		{
			Key:      ConfigDatabaseMaxOpenConns,
			Default:  0,
			Usage:    "Database maximum number of connections in the idle connection pool",
			Type:     config.ValueTypeInt,
			Editable: true,
		},
	}
}

func (r *Resource) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		ConfigDatabaseMigrationsTable: {r.watchMigrationsTable},
		ConfigDatabaseMaxIdleConns:    {r.watchMaxIdleConns},
		ConfigDatabaseMaxOpenConns:    {r.watchMaxOpenConns},
	}
}

func (r *Resource) watchMigrationsTable(newValue interface{}, _ interface{}) {
	migrate.SetTable(newValue.(string))
}

func (r *Resource) watchMaxIdleConns(newValue interface{}, _ interface{}) {
	r.storage.executor.(*gorp.DbMap).Db.SetMaxIdleConns(newValue.(int))
}

func (r *Resource) watchMaxOpenConns(newValue interface{}, _ interface{}) {
	r.storage.executor.(*gorp.DbMap).Db.SetMaxIdleConns(newValue.(int))
}
