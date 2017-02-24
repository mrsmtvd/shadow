package database

import (
	"github.com/kihamo/shadow/components/config"
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

func (c *Component) GetConfigVariables() []config.Variable {
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

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		config.ConfigDebug:            {c.watchDebug},
		ConfigDatabaseMigrationsTable: {c.watchMigrationsTable},
		ConfigDatabaseMaxIdleConns:    {c.watchMaxIdleConns},
		ConfigDatabaseMaxOpenConns:    {c.watchMaxOpenConns},
	}
}

func (c *Component) watchDebug(_ string, newValue interface{}, _ interface{}) {
	c.initTrace(newValue.(bool))
}

func (c *Component) watchMigrationsTable(_ string, newValue interface{}, _ interface{}) {
	migrate.SetTable(newValue.(string))
}

func (c *Component) watchMaxIdleConns(_ string, newValue interface{}, _ interface{}) {
	c.initConns(newValue.(int), c.config.GetInt(ConfigDatabaseMaxOpenConns))
}

func (c *Component) watchMaxOpenConns(_ string, newValue interface{}, _ interface{}) {
	c.initConns(c.config.GetInt(ConfigDatabaseMaxIdleConns), newValue.(int))
}
