package database

import (
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
		config.ConfigDebug:            {r.watchDebug},
		ConfigDatabaseMigrationsTable: {r.watchMigrationsTable},
		ConfigDatabaseMaxIdleConns:    {r.watchMaxIdleConns},
		ConfigDatabaseMaxOpenConns:    {r.watchMaxOpenConns},
	}
}

func (r *Resource) watchDebug(newValue interface{}, _ interface{}) {
	r.initTrace(newValue.(bool))
}

func (r *Resource) watchMigrationsTable(newValue interface{}, _ interface{}) {
	migrate.SetTable(newValue.(string))
}

func (r *Resource) watchMaxIdleConns(newValue interface{}, _ interface{}) {
	r.initConns(newValue.(int), r.config.GetInt(ConfigDatabaseMaxOpenConns))
}

func (r *Resource) watchMaxOpenConns(newValue interface{}, _ interface{}) {
	r.initConns(r.config.GetInt(ConfigDatabaseMaxIdleConns), newValue.(int))
}
