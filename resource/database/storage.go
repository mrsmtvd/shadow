package database

import (
	"database/sql"

	"github.com/dropbox/godropbox/errors"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

func NewSQLStorage(driver string, dataSourceName string) (*SqlStorage, error) {
	db, err := sql.Open(driver, dataSourceName)

	if err != nil {
		return nil, err
	}

	dbMap := &gorp.DbMap{
		Db: db,
	}

	switch driver {
	case "mysql":
		dbMap.Dialect = gorp.MySQLDialect{"InnoDB", "UTF8"}
		break

	default:
		return nil, errors.Newf("Storage driver %s not found", driver)
	}

	return &SqlStorage{
		executor: dbMap,
	}, nil
}

func (s *SqlStorage) CreateTablesIfNotExists() error {
	return s.executor.(*gorp.DbMap).CreateTablesIfNotExists()
}

func (s *SqlStorage) SetTypeConverter(converter gorp.TypeConverter) {
	s.executor.(*gorp.DbMap).TypeConverter = converter
}

func (s *SqlStorage) AddTableWithName(i interface{}, name string) *gorp.TableMap {
	return s.executor.(*gorp.DbMap).AddTableWithName(i, name)
}

func (s *SqlStorage) Begin() (*SqlStorage, error) {
	transaction, err := s.executor.(*gorp.DbMap).Begin()
	if err != nil {
		return nil, err
	}

	return &SqlStorage{
		executor: transaction,
	}, nil
}

func (s *SqlStorage) Commit() error {
	transaction, ok := s.executor.(*gorp.Transaction)
	if !ok {
		return errors.New("Executor is not grop.Transaction")
	}

	return transaction.Commit()
}

func (s *SqlStorage) Rollback() error {
	transaction, ok := s.executor.(*gorp.Transaction)
	if !ok {
		return errors.New("Executor is not grop.Transaction")
	}

	return transaction.Rollback()
}

func (s *SqlStorage) SelectByQuery(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	data, err := s.executor.Select(i, query, args...)
	if err != nil {
		return data, errors.Wrapf(err, "Error getting collection from DB, query: '%s'", query)
	}
	return data, nil
}

func (s *SqlStorage) Select(i interface{}, builder *sq.SelectBuilder) ([]interface{}, error) {
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

func (s *SqlStorage) SelectOne(holder interface{}, builder *sq.SelectBuilder) error {
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

func (s *SqlStorage) Delete(list ...interface{}) (int64, error) {
	count, err := s.executor.Delete(list...)

	if err != nil {
		err = errors.Wrap(err, "Error deleting data in DB")
	}

	return count, err
}

func (s *SqlStorage) Prepare(query string) (*sql.Stmt, error) {
	if transaction, ok := s.executor.(*gorp.Transaction); ok {
		return transaction.Prepare(query)
	} else if dbMap, ok := s.executor.(*gorp.DbMap); ok {
		return dbMap.Prepare(query)
	}

	return nil, errors.New("Executor is not grop.Transaction or gorp.DbMap")
}

func (s *SqlStorage) ExecByQuery(query string, args ...interface{}) (sql.Result, error) {
	result, err := s.executor.Exec(query, args...)
	if err != nil {
		return result, errors.Wrapf(err, "Error executing DB query, query: '%s'", query)
	}
	return result, nil
}

func (s *SqlStorage) Exec(query interface{}, args ...interface{}) (sql.Result, error) {
	if b, ok := query.(*sq.SelectBuilder); ok {
		return s.ExecSelect(b)
	}

	if b, ok := query.(*sq.InsertBuilder); ok {
		return s.ExecInsert(b)
	}

	if b, ok := query.(*sq.UpdateBuilder); ok {
		return s.ExecUpdate(b)
	}

	if b, ok := query.(*sq.DeleteBuilder); ok {
		return s.ExecDelete(b)
	}

	if b, ok := query.(*sq.CaseBuilder); ok {
		return s.ExecCase(b)
	}

	return nil, errors.New("could not prepare SQL query")
}

func (s *SqlStorage) ExecSelect(builder *sq.SelectBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare SQL query")
	}
	return s.ExecByQuery(query, args...)
}

func (s *SqlStorage) ExecInsert(builder *sq.InsertBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare SQL query")
	}
	return s.ExecByQuery(query, args...)
}

func (s *SqlStorage) ExecUpdate(builder *sq.UpdateBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare SQL query")
	}
	return s.ExecByQuery(query, args...)
}

func (s *SqlStorage) ExecDelete(builder *sq.DeleteBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare SQL query")
	}
	return s.ExecByQuery(query, args...)
}

func (s *SqlStorage) ExecCase(builder *sq.CaseBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare SQL query")
	}
	return s.ExecByQuery(query, args...)
}
