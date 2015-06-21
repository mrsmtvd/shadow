package resource

import (
	"github.com/kihamo/shadow/storage"
	"github.com/kihamo/shadow"
)

type Database struct {
	config  *Config
	storage *storage.SqlStorage
}

func (r *Database) GetName() string {
	return "database"
}

func (r *Database) Init(a *shadow.Application) error {
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}

	r.config = resourceConfig.(*Config)
	r.config.Add("database-driver", "mymysql", "Database driver")
	r.config.Add("database-dsn", "tcp:localhost:3306*shadow/root/", "Database DSN")

	return nil
}

func (r *Database) Run() (err error) {
	r.storage, err = storage.NewSQLStorage(r.config.GetString("database-driver"), r.config.GetString("database-dsn"))

	if err != nil {
		return err
	}

	return nil
}

func (r *Database) GetStorage() *storage.SqlStorage {
	return r.storage
}
