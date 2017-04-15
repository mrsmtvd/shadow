package database

import (
	"github.com/kihamo/shadow/components/config"
	"github.com/rubenv/sql-migrate"
)

const (
	ConfigDriver          = ComponentName + ".driver"
	ConfigDsn             = ComponentName + ".dsn"
	ConfigMigrationsTable = ComponentName + ".migrations.table"
	ConfigMaxIdleConns    = ComponentName + ".max_idle_conns"
	ConfigMaxOpenConns    = ComponentName + ".max_open_conns"

	defaultMigrationsTableName = "migrations"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:     ConfigDriver,
			Default: "mysql",
			Usage:   "Database driver (sqlite3, postgres, mysql, mssql and oci8)",
			Type:    config.ValueTypeString,
		},
		{
			Key:   ConfigDsn,
			Usage: "Database DSN",
			Type:  config.ValueTypeString,
		},
		{
			Key:      ConfigMigrationsTable,
			Default:  defaultMigrationsTableName,
			Usage:    "Database migrations table name",
			Type:     config.ValueTypeString,
			Editable: true,
		},
		{
			Key:      ConfigMaxIdleConns,
			Default:  0,
			Usage:    "Database maximum number of connections in the idle connection pool",
			Type:     config.ValueTypeInt,
			Editable: true,
		},
		{
			Key:      ConfigMaxOpenConns,
			Default:  0,
			Usage:    "Database maximum number of connections in the idle connection pool",
			Type:     config.ValueTypeInt,
			Editable: true,
		},
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		config.ConfigDebug:    {c.watchDebug},
		ConfigMigrationsTable: {c.watchMigrationsTable},
		ConfigMaxIdleConns:    {c.watchMaxIdleConns},
		ConfigMaxOpenConns:    {c.watchMaxOpenConns},
	}
}

func (c *Component) watchDebug(_ string, newValue interface{}, _ interface{}) {
	c.initTrace(newValue.(bool))
}

func (c *Component) watchMigrationsTable(_ string, newValue interface{}, _ interface{}) {
	migrate.SetTable(newValue.(string))
}

func (c *Component) watchMaxIdleConns(_ string, newValue interface{}, _ interface{}) {
	c.initConns(newValue.(int), c.config.GetInt(ConfigMaxOpenConns))
}

func (c *Component) watchMaxOpenConns(_ string, newValue interface{}, _ interface{}) {
	c.initConns(c.config.GetInt(ConfigMaxIdleConns), newValue.(int))
}
