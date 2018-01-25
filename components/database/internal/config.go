package internal

import (
	"fmt"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/database"
	"github.com/kihamo/shadow/components/database/storage"
	"github.com/rubenv/sql-migrate"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			database.ConfigAllowUseMasterAsSlave,
			config.ValueTypeBool,
			false,
			"Allow use master as slave",
			true,
			nil,
			nil,
		),
		config.NewVariable(
			database.ConfigBalancer,
			config.ValueTypeString,
			database.BalancerRoundRobin,
			"Balancer for slaves",
			true,
			[]string{
				config.ViewEnum,
			},
			map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{database.BalancerRoundRobin, "Round robin"},
					{database.BalancerRandom, "Random"},
				},
			}),
		config.NewVariable(
			database.ConfigDriver,
			config.ValueTypeString,
			storage.DialectMySQL,
			"Database driver",
			false,
			[]string{
				config.ViewEnum,
			},
			map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{storage.DialectMSSQL, "MSSQL"},
					{storage.DialectMySQL, "MySQL"},
					{storage.DialectOracle, "Oracle"},
					{"oci8", "Oracle"},
					{storage.DialectPostgres, "Postgres"},
					{storage.DialectSQLite3, "SQLite3"},
				},
			}),
		config.NewVariable(
			database.ConfigDialectEngine,
			config.ValueTypeString,
			"InnoDB",
			fmt.Sprintf("Dialect engine (%s)", storage.DialectMySQL),
			false,
			[]string{
				config.ViewEnum,
			},
			map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{storage.EngineInnoDB, storage.EngineInnoDB},
					{storage.EngineMyISAM, storage.EngineMyISAM},
				},
			}),
		config.NewVariable(
			database.ConfigDialectEncoding,
			config.ValueTypeString,
			"UTF8",
			fmt.Sprintf("Dialect encoding (%s)", storage.DialectMySQL),
			false,
			nil,
			nil),
		config.NewVariable(
			database.ConfigDialectVersion,
			config.ValueTypeString,
			nil,
			fmt.Sprintf("Dialect version (%s)", storage.DialectMSSQL),
			false,
			[]string{
				config.ViewEnum,
			},
			map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{storage.Version2005, "Legacy datatypes will be used"},
				},
			}),
		config.NewVariable(
			database.ConfigDsnMaster,
			config.ValueTypeString,
			nil,
			"Database DSN of Master",
			false,
			nil,
			nil),
		config.NewVariable(
			database.ConfigDsnSlaves,
			config.ValueTypeString,
			nil,
			"Database DSN of Slaves",
			false,
			[]string{config.ViewTags},
			map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a slave",
			}),
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
		config.NewWatcher(database.ComponentName, []string{database.ConfigAllowUseMasterAsSlave}, c.watchAllowUseMasterAsSlave),
		config.NewWatcher(database.ComponentName, []string{database.ConfigBalancer}, c.watchBalancer),
		config.NewWatcher(database.ComponentName, []string{config.ConfigDebug}, c.watchDebug),
		config.NewWatcher(database.ComponentName, []string{database.ConfigMigrationsSchema}, c.watchMigrationsSchema),
		config.NewWatcher(database.ComponentName, []string{database.ConfigMigrationsTable}, c.watchMigrationsTable),
		config.NewWatcher(database.ComponentName, []string{database.ConfigMaxIdleConns}, c.watchMaxIdleConns),
		config.NewWatcher(database.ComponentName, []string{database.ConfigMaxOpenConns}, c.watchMaxOpenConns),
	}
}

func (c *Component) watchAllowUseMasterAsSlave(_ string, newValue interface{}, _ interface{}) {
	if newValue.(bool) {
		c.storage.AllowUseMasterAsSlave()
	} else {
		c.storage.DisallowUseMasterAsSlave()
	}
}

func (c *Component) watchBalancer(_ string, _ interface{}, _ interface{}) {
	c.initBalancer()
}

func (c *Component) watchDebug(_ string, newValue interface{}, _ interface{}) {
	c.initTrace()
}

func (c *Component) watchMigrationsTable(_ string, newValue interface{}, _ interface{}) {
	migrate.SetTable(newValue.(string))
}

func (c *Component) watchMigrationsSchema(_ string, newValue interface{}, _ interface{}) {
	migrate.SetSchema(newValue.(string))
}

func (c *Component) watchMaxIdleConns(_ string, newValue interface{}, _ interface{}) {
	c.storage.(*storage.SQL).SetMaxIdleConns(newValue.(int))
}

func (c *Component) watchMaxOpenConns(_ string, newValue interface{}, _ interface{}) {
	c.storage.(*storage.SQL).SetMaxOpenConns(newValue.(int))
}
