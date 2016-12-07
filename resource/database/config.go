package database

import (
	"github.com/kihamo/shadow/resource/config"
)

func (r *Resource) GetConfigVariables() []config.Variable {
	return []config.Variable{
		{
			Key:   "database.driver",
			Value: "mysql",
			Usage: "Database driver (sqlite3, postgres, mysql, mssql and oci8)",
		},
		{
			Key:   "database.dsn",
			Value: "root:@tcp(localhost:3306)/shadow",
			Usage: "Database DSN",
		},
		{
			Key:   "database.migrations.table",
			Value: defaultMigrationsTableName,
			Usage: "Database migrations table name",
		},
		{
			Key:   "database.max_idle_conns",
			Value: 0,
			Usage: "Database maximum number of connections in the idle connection pool",
		},
		{
			Key:   "database.max_open_conns",
			Value: 0,
			Usage: "Database maximum number of connections in the idle connection pool",
		},
	}
}
