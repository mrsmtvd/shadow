package internal

import (
	"fmt"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/database"
	"github.com/rubenv/sql-migrate"
)

const (
	defaultMigrationsTableName = "migrations"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariableItem(
			database.ConfigDriver,
			config.ValueTypeString,
			DialectMySQL,
			fmt.Sprintf("Database driver (%s, %s, %s, %s and %s)", DialectMySQL, DialectOracle, DialectPostgres, DialectSQLite3, DialectMSSQL),
			false,
			[]string{
				config.ViewEnum,
			},
			map[string]interface{}{
				DialectMySQL:    "MySQL",
				DialectOracle:   "Oracle",
				DialectPostgres: "Postgres",
				DialectSQLite3:  "SQLite3",
				DialectMSSQL:    "MSSQL",
			}),
		config.NewVariableItem(
			database.ConfigDialectEngine,
			config.ValueTypeString,
			"InnoDB",
			fmt.Sprintf("Dialect engine (%s)", DialectMySQL),
			false,
			nil,
			nil),
		config.NewVariableItem(
			database.ConfigDialectEncoding,
			config.ValueTypeString,
			"UTF8",
			fmt.Sprintf("Dialect encoding (%s)", DialectMySQL),
			false,
			nil,
			nil),
		config.NewVariableItem(
			database.ConfigDialectVersion,
			config.ValueTypeString,
			nil,
			fmt.Sprintf("Dialect version (%s)", DialectMSSQL),
			false,
			nil,
			nil),
		config.NewVariableItem(
			database.ConfigDsn,
			config.ValueTypeString,
			nil,
			"Database DSN",
			false,
			nil,
			nil),
		config.NewVariableItem(
			database.ConfigMigrationsTable,
			config.ValueTypeString,
			defaultMigrationsTableName,
			"Database migrations table name",
			true,
			nil,
			nil),
		config.NewVariableItem(
			database.ConfigMaxIdleConns,
			config.ValueTypeInt,
			0,
			"Database maximum number of connections in the idle connection pool",
			true,
			nil,
			nil),
		config.NewVariableItem(
			database.ConfigMaxOpenConns,
			config.ValueTypeInt,
			0,
			"Database maximum number of connections in the idle connection pool",
			true,
			nil,
			nil),
	}
}

func (c *Component) GetConfigWatchers() map[string][]config.Watcher {
	return map[string][]config.Watcher{
		config.ConfigDebug:             {c.watchDebug},
		database.ConfigMigrationsTable: {c.watchMigrationsTable},
		database.ConfigMaxIdleConns:    {c.watchMaxIdleConns},
		database.ConfigMaxOpenConns:    {c.watchMaxOpenConns},
	}
}

func (c *Component) watchDebug(_ string, newValue interface{}, _ interface{}) {
	c.initTrace(newValue.(bool))
}

func (c *Component) watchMigrationsTable(_ string, newValue interface{}, _ interface{}) {
	migrate.SetTable(newValue.(string))
}

func (c *Component) watchMaxIdleConns(_ string, newValue interface{}, _ interface{}) {
	c.initConns(newValue.(int), c.config.GetInt(database.ConfigMaxOpenConns))
}

func (c *Component) watchMaxOpenConns(_ string, newValue interface{}, _ interface{}) {
	c.initConns(c.config.GetInt(database.ConfigMaxIdleConns), newValue.(int))
}
