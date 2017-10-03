package database

const (
	ConfigDriver          = ComponentName + ".driver"
	ConfigDialectEngine   = ComponentName + ".dialect.engine"
	ConfigDialectEncoding = ComponentName + ".dialect.encoding"
	ConfigDialectVersion  = ComponentName + ".dialect.version"
	ConfigDsn             = ComponentName + ".dsn"
	ConfigMigrationsTable = ComponentName + ".migrations.table"
	ConfigMaxIdleConns    = ComponentName + ".max_idle_conns"
	ConfigMaxOpenConns    = ComponentName + ".max_open_conns"
)
