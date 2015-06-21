package storage

import (
	"database/sql"

	"github.com/dropbox/godropbox/errors"
	"github.com/lann/squirrel"
	_ "github.com/ziutek/mymysql/godrv"
	"gopkg.in/gorp.v1"
)

type SqlStorage struct {
	executor gorp.SqlExecutor
}

func NewDbMapWrapper(dbMap *gorp.DbMap) *SqlStorage {
	return &SqlStorage{dbMap}
}

func NewSQLStorage(driver string, dataSourceName string) (*SqlStorage, error) {
	db, err := sql.Open(driver, dataSourceName)

	if err != nil {
		return nil, err
	}

	dbMap := &gorp.DbMap{
		Db: db,
	}

	switch driver {
	case "mymysql":
		dbMap.Dialect = gorp.MySQLDialect{"InnoDB", "UTF8"}
		break

	default:
		return nil, errors.Newf("Storage driver %s not found", driver)
	}

	return NewDbMapWrapper(dbMap), nil
}

func (s *SqlStorage) CreateTablesIfNotExists() error {
	return s.executor.(*gorp.DbMap).CreateTablesIfNotExists()
}

func (s *SqlStorage) AddTableWithName(i interface{}, name string) *gorp.TableMap {
	return s.executor.(*gorp.DbMap).AddTableWithName(i, name)
}

func (s *SqlStorage) SelectByQuery(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	data, err := s.executor.Select(i, query, args...)
	if err != nil {
		return data, errors.Wrapf(err, "Error getting collection from DB, query: '%s'", query)
	}
	return data, nil
}

func (s *SqlStorage) Select(i interface{}, builder *squirrel.SelectBuilder) ([]interface{}, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare SQL query")
	}
	return s.SelectByQuery(i, query, args...)
}

func (s *SqlStorage) SelectOneByQuery(holder interface{}, query string, args ...interface{}) error {
	err := s.executor.SelectOne(holder, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrapf(err, "Error getting value from DB, query: '%s'", query)
	}
	return err
}

func (s *SqlStorage) SelectOne(holder interface{}, builder *squirrel.SelectBuilder) error {
	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "could not prepare SQL query")
	}
	return s.SelectOneByQuery(holder, query, args...)
}

func (s *SqlStorage) Insert(list ...interface{}) error {
	if err := s.executor.Insert(list...); err != nil {
		return errors.Wrap(err, "Error inserting data into DB")
	}

	return nil
}

func (s *SqlStorage) Update(list ...interface{}) (int64, error) {
	count, err := s.executor.Update(list...)

	if err != nil {
		err = errors.Wrap(err, "Error updating data in DB")
	}

	return count, err
}
