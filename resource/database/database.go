package database

import (
	"github.com/go-gorp/gorp"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
)

type SqlStorage struct {
	executor gorp.SqlExecutor
}

type Database struct {
	application *shadow.Application
	config      *resource.Config
	storage     *SqlStorage
}

func (r *Database) GetName() string {
	return "database"
}

func (r *Database) GetConfigVariables() []resource.ConfigVariable {
	return []resource.ConfigVariable{
		resource.ConfigVariable{
			Key:   "database.driver",
			Value: "mysql",
			Usage: "Database driver",
		},
		resource.ConfigVariable{
			Key:   "database.dsn",
			Value: "root:@tcp(localhost:3306)/shadow",
			Usage: "Database DSN",
		},
	}
}

func (r *Database) Init(a *shadow.Application) error {
	r.application = a
	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}

	r.config = resourceConfig.(*resource.Config)

	return nil
}

func (r *Database) Run() (err error) {
	r.storage, err = NewSQLStorage(r.config.GetString("database.driver"), r.config.GetString("database.dsn"))

	if err != nil {
		return err
	}

	if r.config.GetBool("debug") {
		resourceLogger, err := r.application.GetResource("logger")
		if err != nil {
			return err
		}

		logger := resourceLogger.(*resource.Logger).Get(r.GetName())
		r.storage.executor.(*gorp.DbMap).TraceOn("", logger)
	}

	r.storage.SetTypeConverter(TypeConverter{})

	return nil
}

func (r *Database) GetStorage() *SqlStorage {
	return r.storage
}
