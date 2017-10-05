package internal

import (
	"fmt"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/database"
	"github.com/rubenv/sql-migrate"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			database.ConfigDriver,
			config.ValueTypeString,
			DialectMySQL,
			"Database driver",
			false,
			[]string{
				config.ViewEnum,
			},
			map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{DialectMSSQL, "MSSQL"},
					{DialectMySQL, "MySQL"},
					{DialectOracle, "Oracle"},
					{"oci8", "Oracle"},
					{DialectPostgres, "Postgres"},
					{DialectSQLite3, "SQLite3"},
				},
			}),
		config.NewVariable(
			database.ConfigDialectEngine,
			config.ValueTypeString,
			"InnoDB",
			fmt.Sprintf("Dialect engine (%s)", DialectMySQL),
			false,
			[]string{
				config.ViewEnum,
			},
			map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{EngineInnoDB, EngineInnoDB},
					{EngineMyISAM, EngineMyISAM},
				},
			}),
		config.NewVariable(
			database.ConfigDialectEncoding,
			config.ValueTypeString,
			"UTF8",
			fmt.Sprintf("Dialect encoding (%s)", DialectMySQL),
			false,
			nil,
			nil),
		config.NewVariable(
			database.ConfigDialectVersion,
			config.ValueTypeString,
			nil,
			fmt.Sprintf("Dialect version (%s)", DialectMSSQL),
			false,
			[]string{
				config.ViewEnum,
			},
			map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{Version2005, "Legacy datatypes will be used"},
				},
			}),
		config.NewVariable(
			database.ConfigDsn,
			config.ValueTypeString,
			nil,
			"Database DSN",
			false,
			nil,
			nil),
		config.NewVariable(
			database.ConfigMigrationsSchema,
			config.ValueTypeString,
			"",
			"Database migrations schema name",
			true,
			nil,
			nil),
		config.NewVariable(
			database.ConfigMigrationsTable,
			config.ValueTypeString,
			"migrations",
			"Database migrations table name",
			true,
			nil,
			nil),
		config.NewVariable(
			database.ConfigMaxIdleConns,
			config.ValueTypeInt,
			0,
			"Database maximum number of connections in the idle connection pool",
			true,
			nil,
			nil),
		config.NewVariable(
			database.ConfigMaxOpenConns,
			config.ValueTypeInt,
			0,
			"Database maximum number of connections in the idle connection pool",
			true,
			nil,
			nil),
	}
}

func (c *Component) GetConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher(database.ComponentName, []string{config.ConfigDebug}, c.watchDebug),
		config.NewWatcher(database.ComponentName, []string{database.ConfigMigrationsSchema}, c.watchMigrationsSchema),
		config.NewWatcher(database.ComponentName, []string{database.ConfigMigrationsTable}, c.watchMigrationsTable),
		config.NewWatcher(database.ComponentName, []string{
			database.ConfigMaxIdleConns,
			database.ConfigMaxOpenConns,
		}, c.watchFoxMaxConns),
	}
}

func (c *Component) watchDebug(_ string, newValue interface{}, _ interface{}) {
	c.initTrace(newValue.(bool))
}

func (c *Component) watchMigrationsTable(_ string, newValue interface{}, _ interface{}) {
	migrate.SetTable(newValue.(string))
}

func (c *Component) watchMigrationsSchema(_ string, newValue interface{}, _ interface{}) {
	migrate.SetSchema(newValue.(string))
}

func (c *Component) watchFoxMaxConns(_ string, _ interface{}, _ interface{}) {
	c.initConns(c.config.GetInt(database.ConfigMaxIdleConns), c.config.GetInt(database.ConfigMaxOpenConns))
}
