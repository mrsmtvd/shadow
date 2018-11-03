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
		config.NewVariable(database.ConfigDriver, config.ValueTypeString).
			WithUsage("Driver").
			WithGroup("Driver").
			WithDefault(storage.DialectMySQL).
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{storage.DialectMSSQL, "MSSQL"},
					{storage.DialectMySQL, "MySQL"},
					{storage.DialectOracle, "Oracle"},
					{"oci8", "Oracle"},
					{storage.DialectPostgres, "Postgres"},
					{storage.DialectSQLite3, "SQLite3"},
				},
			}),
		config.NewVariable(database.ConfigDialectEngine, config.ValueTypeString).
			WithUsage(fmt.Sprintf("Dialect engine for %s driver", storage.DialectMySQL)).
			WithGroup("Driver").
			WithDefault(storage.EngineInnoDB).
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{storage.EngineInnoDB, storage.EngineInnoDB},
					{storage.EngineMyISAM, storage.EngineMyISAM},
				},
			}),
		config.NewVariable(database.ConfigDialectEncoding, config.ValueTypeString).
			WithUsage(fmt.Sprintf("Dialect encoding for %s driver", storage.DialectMySQL)).
			WithGroup("Driver").
			WithDefault("UTF8"),
		config.NewVariable(database.ConfigDialectVersion, config.ValueTypeString).
			WithUsage(fmt.Sprintf("Dialect version for %s driver", storage.DialectMSSQL)).
			WithGroup("Driver").
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{storage.Version2005, "Legacy datatypes will be used"},
				},
			}),
		config.NewVariable(database.ConfigDsnMaster, config.ValueTypeString).
			WithUsage("DSN of master").
			WithGroup("Master-Slave"),
		config.NewVariable(database.ConfigDsnSlaves, config.ValueTypeString).
			WithUsage("DSN of slaves").
			WithGroup("Master-Slave"),
		config.NewVariable(database.ConfigAllowUseMasterAsSlave, config.ValueTypeBool).
			WithUsage("Allow use master as slave").
			WithGroup("Master-Slave").
			WithEditable(true),
		config.NewVariable(database.ConfigBalancer, config.ValueTypeString).
			WithUsage("Balancer for slaves").
			WithGroup("Master-Slave").
			WithEditable(true).
			WithDefault(database.BalancerRoundRobin).
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{database.BalancerRoundRobin, "Round robin"},
					{database.BalancerRandom, "Random"},
				},
			}),
		config.NewVariable(database.ConfigMigrationsSchema, config.ValueTypeString).
			WithUsage("Migrations schema name").
			WithGroup("Migrations").
			WithEditable(true),
		config.NewVariable(database.ConfigMigrationsTable, config.ValueTypeString).
			WithUsage("Migrations table name").
			WithGroup("Migrations").
			WithEditable(true).
			WithDefault("migrations"),
		config.NewVariable(database.ConfigMaxIdleConns, config.ValueTypeInt).
			WithUsage("Maximum number of connections in the idle connection pool").
			WithGroup("Connections").
			WithEditable(true).
			WithDefault(0),
		config.NewVariable(database.ConfigMaxOpenConns, config.ValueTypeInt).
			WithUsage("Maximum number of connections in the open connection pool").
			WithGroup("Connections").
			WithEditable(true).
			WithDefault(0),
		config.NewVariable(database.ConfigConnMaxLifetime, config.ValueTypeDuration).
			WithUsage("Maximum amount of time a connection may be reused").
			WithGroup("Connections").
			WithEditable(true).
			WithDefault(0),
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
	if s := c.Storage(); s != nil {
		if newValue.(bool) {
			c.Storage().AllowUseMasterAsSlave()
		} else {
			c.Storage().DisallowUseMasterAsSlave()
		}
	}
}

func (c *Component) watchBalancer(_ string, newValue interface{}, _ interface{}) {
	if s := c.Storage(); s != nil {
		c.initBalancer(s, newValue.(string))
	}
}

func (c *Component) watchDebug(_ string, newValue interface{}, _ interface{}) {
	if s := c.Storage(); s != nil {
		c.initTrace(s, newValue.(bool))
	}
}

func (c *Component) watchMigrationsTable(_ string, newValue interface{}, _ interface{}) {
	migrate.SetTable(newValue.(string))
}

func (c *Component) watchMigrationsSchema(_ string, newValue interface{}, _ interface{}) {
	migrate.SetSchema(newValue.(string))
}

func (c *Component) watchMaxIdleConns(_ string, newValue interface{}, _ interface{}) {
	if s := c.Storage(); s != nil {
		s.(*storage.SQL).SetMaxIdleConns(newValue.(int))
	}
}

func (c *Component) watchMaxOpenConns(_ string, newValue interface{}, _ interface{}) {
	if s := c.Storage(); s != nil {
		s.(*storage.SQL).SetMaxOpenConns(newValue.(int))
	}
}

func (c *Component) watchConnMaxLifetime(_ string, newValue interface{}, _ interface{}) {
	if s := c.Storage(); s != nil {
		s.(*storage.SQL).SetConnMaxLifetime(newValue.(time.Duration))
	}
}
