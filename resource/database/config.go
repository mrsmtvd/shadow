package database

import (
	"github.com/kihamo/shadow/resource/config"
)

func (r *Resource) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:     "database.driver",
			Default: "mysql",
			Usage:   "Database driver (sqlite3, postgres, mysql, mssql and oci8)",
			Type:    config.ValueTypeString,
		},
		{
			Key:   "database.dsn",
			Usage: "Database DSN",
			Type:  config.ValueTypeString,
		},
		{
			Key:     "database.migrations.table",
			Default: defaultMigrationsTableName,
			Usage:   "Database migrations table name",
			Type:    config.ValueTypeString,
		},
		{
			Key:     "database.max_idle_conns",
			Default: 0,
			Usage:   "Database maximum number of connections in the idle connection pool",
			Type:    config.ValueTypeInt,
		},
		{
			Key:     "database.max_open_conns",
			Default: 0,
			Usage:   "Database maximum number of connections in the idle connection pool",
			Type:    config.ValueTypeInt,
		},
	}
}
