package database

const (
	ConfigAllowUseMasterAsSlave = ComponentName + ".allow-use-master-as-slave"
	ConfigBalancer              = ComponentName + ".balancer"
	ConfigDriver                = ComponentName + ".driver"
	ConfigDialectEngine         = ComponentName + ".dialect.engine"
	ConfigDialectEncoding       = ComponentName + ".dialect.encoding"
	ConfigDialectVersion        = ComponentName + ".dialect.version"
	ConfigDsnMaster             = ComponentName + ".dsn.master"
	ConfigDsnSlaves             = ComponentName + ".dsn.slaves"
	ConfigMigrationsSchema      = ComponentName + ".migrations.schema"
	ConfigMigrationsTable       = ComponentName + ".migrations.table"
	ConfigMaxIdleConns          = ComponentName + ".max_idle_conns"
	ConfigMaxOpenConns          = ComponentName + ".max_open_conns"
	ConfigConnMaxLifetime       = ComponentName + ".conn_max_lifetime"
)
