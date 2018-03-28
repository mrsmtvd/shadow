package internal

import (
	"fmt"
	"time"

	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/database"
	"github.com/kihamo/shadow/components/database/storage"
	"github.com/rubenv/sql-migrate"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			database.ConfigDriver,
			config.ValueTypeString,
			storage.DialectMySQL,
			"Driver",
			false,
			"Driver",
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
			fmt.Sprintf("Dialect engine for %s driver", storage.DialectMySQL),
			false,
			"Driver",
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
			fmt.Sprintf("Dialect encoding for %s driver", storage.DialectMySQL),
			false,
			"Driver",
			nil,
			nil),
		config.NewVariable(
			database.ConfigDialectVersion,
			config.ValueTypeString,
			nil,
			fmt.Sprintf("Dialect version for %s driver", storage.DialectMSSQL),
			false,
			"Driver",
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
			"DSN of master",
			false,
			"Master-Slave",
			nil,
			nil),
		config.NewVariable(
			database.ConfigDsnSlaves,
			config.ValueTypeString,
			nil,
			"DSN of slaves",
			false,
			"Master-Slave",
			[]string{config.ViewTags},
			map[string]interface{}{
				config.ViewOptionTagsDefaultText: "add a slave",
			}),
		config.NewVariable(
			database.ConfigAllowUseMasterAsSlave,
			config.ValueTypeBool,
			false,
			"Allow use master as slave",
			true,
			"",
			nil,
			nil,
		),
		config.NewVariable(
			database.ConfigBalancer,
			config.ValueTypeString,
			database.BalancerRoundRobin,
			"Balancer for slaves",
			true,
			"",
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
			database.ConfigMigrationsSchema,
			config.ValueTypeString,
			"",
			"Migrations schema name",
			true,
			"Migrations",
			nil,
			nil),
		config.NewVariable(
			database.ConfigMigrationsTable,
			config.ValueTypeString,
			"migrations",
			"Migrations table name",
			true,
			"Migrations",
			nil,
			nil),
		config.NewVariable(
			database.ConfigMaxIdleConns,
			config.ValueTypeInt,
			0,
			"Maximum number of connections in the idle connection pool",
			true,
			"Connections",
			nil,
			nil),
		config.NewVariable(
			database.ConfigMaxOpenConns,
			config.ValueTypeInt,
			0,
			"Maximum number of connections in the open connection pool",
			true,
			"Connections",
			nil,
			nil),
		config.NewVariable(
			database.ConfigConnMaxLifetime,
			config.ValueTypeDuration,
			0,
			"Maximum amount of time a connection may be reused",
			true,
			"Connections",
			nil,
			nil),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{database.ConfigAllowUseMasterAsSlave}, c.watchAllowUseMasterAsSlave),
		config.NewWatcher([]string{database.ConfigBalancer}, c.watchBalancer),
		config.NewWatcher([]string{config.ConfigDebug}, c.watchDebug),
		config.NewWatcher([]string{database.ConfigMigrationsSchema}, c.watchMigrationsSchema),
		config.NewWatcher([]string{database.ConfigMigrationsTable}, c.watchMigrationsTable),
		config.NewWatcher([]string{database.ConfigMaxIdleConns}, c.watchMaxIdleConns),
		config.NewWatcher([]string{database.ConfigMaxOpenConns}, c.watchMaxOpenConns),
		config.NewWatcher([]string{database.ConfigConnMaxLifetime}, c.watchConnMaxLifetime),
	}
}

func (c *Component) watchAllowUseMasterAsSlave(_ string, newValue interface{}, _ interface{}) {
	if newValue.(bool) {
		c.Storage().AllowUseMasterAsSlave()
	} else {
		c.Storage().DisallowUseMasterAsSlave()
	}
}

func (c *Component) watchBalancer(_ string, newValue interface{}, _ interface{}) {
	c.initBalancer(c.Storage(), newValue.(string))
}

func (c *Component) watchDebug(_ string, newValue interface{}, _ interface{}) {
	c.initTrace(c.Storage(), newValue.(bool))
}

func (c *Component) watchMigrationsTable(_ string, newValue interface{}, _ interface{}) {
	migrate.SetTable(newValue.(string))
}

func (c *Component) watchMigrationsSchema(_ string, newValue interface{}, _ interface{}) {
	migrate.SetSchema(newValue.(string))
}

func (c *Component) watchMaxIdleConns(_ string, newValue interface{}, _ interface{}) {
	c.Storage().(*storage.SQL).SetMaxIdleConns(newValue.(int))
}

func (c *Component) watchMaxOpenConns(_ string, newValue interface{}, _ interface{}) {
	c.Storage().(*storage.SQL).SetMaxOpenConns(newValue.(int))
}

func (c *Component) watchConnMaxLifetime(_ string, newValue interface{}, _ interface{}) {
	c.Storage().(*storage.SQL).SetConnMaxLifetime(newValue.(time.Duration))
}
