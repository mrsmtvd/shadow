package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-gorp/gorp"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

const (
	DialectMySQL    = "mysql"
	DialectOracle   = "oracle"
	DialectPostgres = "postgres"
	DialectSQLite3  = "sqlite3"
	DialectMSSQL    = "mssql"

	DialectOptionEngine   = "engine"
	DialectOptionEncoding = "encoding"
	DialectOptionVersion  = "version"
)

type SqlStorage struct {
	executor gorp.SqlExecutor
	dialect  string
}

func NewSQLStorage(driver string, dataSourceName string, options map[string]string) (*SqlStorage, error) {
	var dialect gorp.Dialect
	dialectName := driver

	switch driver {
	case DialectMySQL:
		var (
			ok       bool
			engine   string
			encoding string
		)

		if engine, ok = options[DialectOptionEngine]; !ok {
			engine = "InnoDB"
		}

		if encoding, ok = options[DialectOptionEncoding]; !ok {
			encoding = "UTF8"
		}

		dialect = gorp.MySQLDialect{
			Engine:   engine,
			Encoding: encoding,
		}
		break

	case DialectOracle, "oci8":
		dialect = gorp.OracleDialect{}
		driver = "oci8"
		dialectName = DialectOracle
		break

	case DialectPostgres:
		dialect = gorp.PostgresDialect{}
		break

	case DialectSQLite3:
		dialect = gorp.SqliteDialect{}
		break

	case DialectMSSQL:
		version := ""

		if v, ok := options[DialectOptionVersion]; ok {
			version = v
		}

		dialect = gorp.SqlServerDialect{
			Version: version,
		}
		break

	default:
		return nil, fmt.Errorf("Storage driver %s not found", driver)
	}

	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	dbMap := &gorp.DbMap{
		Db:      db,
		Dialect: dialect,
	}

	return &SqlStorage{
		executor: dbMap,
		dialect:  dialectName,
	}, nil
}

func (s *SqlStorage) GetDialect() string {
	return s.dialect
}

func (s *SqlStorage) CreateTablesIfNotExists() error {
	defer updateMetric(OperationCreate, time.Now())

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
	defer updateMetric(OperationSelect, time.Now())

	data, err := s.executor.Select(i, query, args...)
	if err != nil {
		return data, fmt.Errorf("Error getting collection from DB, query: '%s', error: '%s'", query, err.Error())
	}
	return data, nil
}

func (s *SqlStorage) Select(i interface{}, builder *sq.SelectBuilder) ([]interface{}, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.New("could not prepare SQL query, error: '" + err.Error() + "'")
	}
	return s.SelectByQuery(i, query, args...)
}

func (s *SqlStorage) SelectOneByQuery(holder interface{}, query string, args ...interface{}) error {
	defer updateMetric(OperationSelect, time.Now())

	err := s.executor.SelectOne(holder, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("Error getting value from DB, query: '%s', error: '%s'", query, err.Error())
	}
	return err
}

func (s *SqlStorage) SelectOne(holder interface{}, builder *sq.SelectBuilder) error {
	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("could not prepare SQL query, error: '%s'", err.Error())
	}
	return s.SelectOneByQuery(holder, query, args...)
}

func (s *SqlStorage) SelectIntByQuery(query string, args ...interface{}) (int64, error) {
	defer updateMetric(OperationSelect, time.Now())

	result, err := s.executor.SelectInt(query, args...)

	if err != nil {
		err = fmt.Errorf("Error selecting data in DB, error: '%s'", err.Error())
	}

	return result, err
}

func (s *SqlStorage) SelectInt(builder *sq.SelectBuilder) (int64, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return -1, fmt.Errorf("could not prepare SQL query, error: '%s'", err.Error())
	}
	return s.SelectIntByQuery(query, args...)
}

func (s *SqlStorage) SelectNullIntByQuery(query string, args ...interface{}) (sql.NullInt64, error) {
	defer updateMetric(OperationSelect, time.Now())

	result, err := s.executor.SelectNullInt(query, args...)

	if err != nil {
		err = fmt.Errorf("Error selecting data in DB, error: '%s'", err.Error())
	}

	return result, err
}

func (s *SqlStorage) SelectNullInt(builder *sq.SelectBuilder) (sql.NullInt64, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		var h sql.NullInt64

		return h, fmt.Errorf("could not prepare SQL query, error: '%s'", err.Error())
	}
	return s.SelectNullIntByQuery(query, args...)
}

func (s *SqlStorage) SelectFloatByQuery(query string, args ...interface{}) (float64, error) {
	defer updateMetric(OperationSelect, time.Now())

	result, err := s.executor.SelectFloat(query, args...)

	if err != nil {
		err = fmt.Errorf("Error selecting data in DB, error: '%s'", err.Error())
	}

	return result, err
}

func (s *SqlStorage) SelectFloat(builder *sq.SelectBuilder) (float64, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return -1, fmt.Errorf("could not prepare SQL query, error: '%s'", err.Error())
	}
	return s.SelectFloatByQuery(query, args...)
}

func (s *SqlStorage) SelectNullFloatByQuery(query string, args ...interface{}) (sql.NullFloat64, error) {
	defer updateMetric(OperationSelect, time.Now())

	result, err := s.executor.SelectNullFloat(query, args...)

	if err != nil {
		err = fmt.Errorf("Error selecting data in DB, error: '%s'", err.Error())
	}

	return result, err
}

func (s *SqlStorage) SelectNullFloat(builder *sq.SelectBuilder) (sql.NullFloat64, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		var h sql.NullFloat64

		return h, fmt.Errorf("could not prepare SQL query, error: '%s'", err.Error())
	}
	return s.SelectNullFloatByQuery(query, args...)
}

func (s *SqlStorage) SelectStrByQuery(query string, args ...interface{}) (string, error) {
	defer updateMetric(OperationSelect, time.Now())

	result, err := s.executor.SelectStr(query, args...)

	if err != nil {
		err = fmt.Errorf("Error selecting data in DB, error: '%s'", err.Error())
	}

	return result, err
}

func (s *SqlStorage) SelectStr(builder *sq.SelectBuilder) (string, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return "", fmt.Errorf("could not prepare SQL query, error: '%s'", err.Error())
	}
	return s.SelectStrByQuery(query, args...)
}

func (s *SqlStorage) SelectNullStrByQuery(query string, args ...interface{}) (sql.NullString, error) {
	defer updateMetric(OperationSelect, time.Now())

	result, err := s.executor.SelectNullStr(query, args...)

	if err != nil {
		err = fmt.Errorf("Error selecting data in DB, error: '%s'", err.Error())
	}

	return result, err
}

func (s *SqlStorage) SelectNullStr(builder *sq.SelectBuilder) (sql.NullString, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		var h sql.NullString

		return h, fmt.Errorf("could not prepare SQL query, error: '%s'", err.Error())
	}
	return s.SelectNullStrByQuery(query, args...)
}

func (s *SqlStorage) Get(i interface{}, keys ...interface{}) (interface{}, error) {
	defer updateMetric(OperationSelect, time.Now())

	entity, err := s.executor.Get(i, keys...)

	if err != nil {
		err = fmt.Errorf("Error get data in DB, error: '%s'", err.Error())
	}

	return entity, err
}

func (s *SqlStorage) Insert(list ...interface{}) error {
	defer updateMetric(OperationInsert, time.Now())

	if err := s.executor.Insert(list...); err != nil {
		return fmt.Errorf("Error inserting data into DB, error: '%s'", err.Error())
	}

	return nil
}

func (s *SqlStorage) Update(list ...interface{}) (int64, error) {
	defer updateMetric(OperationUpdate, time.Now())

	count, err := s.executor.Update(list...)

	if err != nil {
		err = fmt.Errorf("Error updating data in DB, error: '%s'", err.Error())
	}

	return count, err
}

func (s *SqlStorage) Delete(list ...interface{}) (int64, error) {
	defer updateMetric(OperationDelete, time.Now())

	count, err := s.executor.Delete(list...)

	if err != nil {
		err = fmt.Errorf("Error deleting data in DB, error: '%s'", err.Error())
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
	defer updateMetric(OperationExec, time.Now())

	result, err := s.executor.Exec(query, args...)
	if err != nil {
		return result, fmt.Errorf("Error executing DB query, query: '%s', error: '%s'", query, err.Error())
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
		return nil, fmt.Errorf("could not prepare SQL query, error: '%s'", err.Error())
	}
	return s.ExecByQuery(query, args...)
}

func (s *SqlStorage) ExecInsert(builder *sq.InsertBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("could not prepare SQL query, error: '%s'", err.Error())
	}
	return s.ExecByQuery(query, args...)
}

func (s *SqlStorage) ExecUpdate(builder *sq.UpdateBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("could not prepare SQL query, error: '%s'", err.Error())
	}
	return s.ExecByQuery(query, args...)
}

func (s *SqlStorage) ExecDelete(builder *sq.DeleteBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("could not prepare SQL query, error: '%s'", err.Error())
	}
	return s.ExecByQuery(query, args...)
}

func (s *SqlStorage) ExecCase(builder *sq.CaseBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("could not prepare SQL query, error: '%s'", err.Error())
	}
	return s.ExecByQuery(query, args...)
}
