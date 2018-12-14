package storage

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/kihamo/shadow/components/database"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

const (
	DialectMySQL    = "mysql"
	DialectOracle   = "oracle"
	DialectPostgres = "postgres"
	DialectSQLite3  = "sqlite3"
	DialectMSSQL    = "mssql"

	EngineInnoDB = "InnoDB"
	EngineMyISAM = "MyISAM"

	Version2005 = "2005"

	DialectOptionEngine   = "engine"
	DialectOptionEncoding = "encoding"
	DialectOptionVersion  = "version"
)

var dsnPattern = regexp.MustCompile(
	`^(?:(?P<user>.*?)(?::(?P<passwd>.*))?@)?` + // [user[:password]@]
		`(?:(?P<net>[^\(]*)(?:\((?P<addr>[^\)]*)\))?)?` + // [net[(addr)]]
		`\/(?P<dbname>.*?)` + // /dbname
		`(?:\?(?P<params>[^\?]*))?$`) // [?param1=value1&paramN=valueN]

type SQLExecutor struct {
	executor      gorp.SqlExecutor
	dialect       string
	name          string
	serverAddress string
}

func NewSQLExecutor(driver string, dataSourceName string, options map[string]string) (*SQLExecutor, error) {
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

	case DialectOracle, "oci8":
		dialect = gorp.OracleDialect{}
		driver = "oci8"
		dialectName = DialectOracle

	case DialectPostgres:
		dialect = gorp.PostgresDialect{}

	case DialectSQLite3:
		dialect = gorp.SqliteDialect{}

	case DialectMSSQL:
		version := ""

		if v, ok := options[DialectOptionVersion]; ok {
			version = v
		}

		dialect = gorp.SqlServerDialect{
			Version: version,
		}

	default:
		return nil, errors.New("executor driver " + driver + " not found")
	}

	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	dbMap := &gorp.DbMap{
		Db:      db,
		Dialect: dialect,
	}

	e := &SQLExecutor{
		executor: dbMap,
		dialect:  dialectName,
		name:     dialectName + ">",
	}

	matches := dsnPattern.FindStringSubmatch(dataSourceName)
	names := dsnPattern.SubexpNames()

	for i, match := range matches {
		switch names[i] {
		case "user":
			e.name += match + ":"
		case "passwd":
			e.name += "*****"
		case "net":
			e.name += "@" + match
		case "addr":
			e.name += "(" + match + ")"
			e.serverAddress = match
		case "dbname":
			e.name += "/" + match
		}
	}

	return e, nil
}

func (e *SQLExecutor) String() string {
	return e.name
}

func (e *SQLExecutor) ServerAddress() string {
	return e.serverAddress
}

func (e *SQLExecutor) DB() *sql.DB {
	return e.executor.(*gorp.DbMap).Db
}

func (e *SQLExecutor) Ping(ctx context.Context) error {
	return e.DB().PingContext(ctx)
}

func (e *SQLExecutor) Begin() (database.Executor, error) {
	transaction, err := e.executor.(*gorp.DbMap).Begin()
	if err != nil {
		return nil, err
	}

	return &SQLExecutor{
		executor: transaction,
	}, nil
}

func (e *SQLExecutor) Commit() error {
	transaction, ok := e.executor.(*gorp.Transaction)
	if !ok {
		return errors.New("executor is not grop.Transaction")
	}

	return transaction.Commit()
}

func (e *SQLExecutor) Rollback() error {
	transaction, ok := e.executor.(*gorp.Transaction)
	if !ok {
		return errors.New("executor is not grop.Transaction")
	}

	return transaction.Rollback()
}

func (e *SQLExecutor) SelectByQuery(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	defer UpdateStorageSQLMetric(OperationSelect, e.serverAddress, time.Now())

	data, err := e.executor.Select(i, query, args...)
	if err != nil {
		return data, errors.New("error getting collection from DB, query: '" + query + "', error: '" + err.Error() + "'")
	}
	return data, nil
}

func (e *SQLExecutor) Select(i interface{}, builder *sq.SelectBuilder) ([]interface{}, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.New("could not prepare SQL query, error: '" + err.Error() + "'")
	}
	return e.SelectByQuery(i, query, args...)
}

func (e *SQLExecutor) SelectOneByQuery(holder interface{}, query string, args ...interface{}) error {
	defer UpdateStorageSQLMetric(OperationSelect, e.serverAddress, time.Now())

	err := e.executor.SelectOne(holder, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return errors.New("error getting value from DB, query: '" + query + "', error: '" + err.Error() + "'")
	}
	return err
}

func (e *SQLExecutor) SelectOne(holder interface{}, builder *sq.SelectBuilder) error {
	query, args, err := builder.ToSql()
	if err != nil {
		return errors.New("could not prepare SQL query, error: '" + err.Error() + "'")
	}
	return e.SelectOneByQuery(holder, query, args...)
}

func (e *SQLExecutor) SelectIntByQuery(query string, args ...interface{}) (int64, error) {
	defer UpdateStorageSQLMetric(OperationSelect, e.serverAddress, time.Now())

	result, err := e.executor.SelectInt(query, args...)

	if err != nil {
		err = errors.New("error selecting data in DB, error: '" + err.Error() + "'")
	}

	return result, err
}

func (e *SQLExecutor) SelectInt(builder *sq.SelectBuilder) (int64, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return -1, errors.New("could not prepare SQL query, error: '" + err.Error() + "'")
	}
	return e.SelectIntByQuery(query, args...)
}

func (e *SQLExecutor) SelectNullIntByQuery(query string, args ...interface{}) (sql.NullInt64, error) {
	defer UpdateStorageSQLMetric(OperationSelect, e.serverAddress, time.Now())

	result, err := e.executor.SelectNullInt(query, args...)

	if err != nil {
		err = errors.New("error selecting data in DB, error: '" + err.Error() + "'")
	}

	return result, err
}

func (e *SQLExecutor) SelectNullInt(builder *sq.SelectBuilder) (sql.NullInt64, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		var h sql.NullInt64

		return h, errors.New("could not prepare SQL query, error: '" + err.Error() + "'")
	}
	return e.SelectNullIntByQuery(query, args...)
}

func (e *SQLExecutor) SelectFloatByQuery(query string, args ...interface{}) (float64, error) {
	defer UpdateStorageSQLMetric(OperationSelect, e.serverAddress, time.Now())

	result, err := e.executor.SelectFloat(query, args...)

	if err != nil {
		err = errors.New("error selecting data in DB, error: '" + err.Error() + "'")
	}

	return result, err
}

func (e *SQLExecutor) SelectFloat(builder *sq.SelectBuilder) (float64, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return -1, errors.New("could not prepare SQL query, error: '" + err.Error() + "'")
	}
	return e.SelectFloatByQuery(query, args...)
}

func (e *SQLExecutor) SelectNullFloatByQuery(query string, args ...interface{}) (sql.NullFloat64, error) {
	defer UpdateStorageSQLMetric(OperationSelect, e.serverAddress, time.Now())

	result, err := e.executor.SelectNullFloat(query, args...)

	if err != nil {
		err = errors.New("error selecting data in DB, error: '" + err.Error() + "'")
	}

	return result, err
}

func (e *SQLExecutor) SelectNullFloat(builder *sq.SelectBuilder) (sql.NullFloat64, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		var h sql.NullFloat64

		return h, errors.New("could not prepare SQL query, error: '" + err.Error() + "'")
	}
	return e.SelectNullFloatByQuery(query, args...)
}

func (e *SQLExecutor) SelectStrByQuery(query string, args ...interface{}) (string, error) {
	defer UpdateStorageSQLMetric(OperationSelect, e.serverAddress, time.Now())

	result, err := e.executor.SelectStr(query, args...)

	if err != nil {
		err = errors.New("error selecting data in DB, error: '" + err.Error() + "'")
	}

	return result, err
}

func (e *SQLExecutor) SelectStr(builder *sq.SelectBuilder) (string, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return "", errors.New("could not prepare SQL query, error: '" + err.Error() + "'")
	}
	return e.SelectStrByQuery(query, args...)
}

func (e *SQLExecutor) SelectNullStrByQuery(query string, args ...interface{}) (sql.NullString, error) {
	defer UpdateStorageSQLMetric(OperationSelect, e.serverAddress, time.Now())

	result, err := e.executor.SelectNullStr(query, args...)

	if err != nil {
		err = errors.New("error selecting data in DB, error: '" + err.Error() + "'")
	}

	return result, err
}

func (e *SQLExecutor) SelectNullStr(builder *sq.SelectBuilder) (sql.NullString, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		var h sql.NullString

		return h, errors.New("could not prepare SQL query, error: '" + err.Error() + "'")
	}
	return e.SelectNullStrByQuery(query, args...)
}

func (e *SQLExecutor) Get(i interface{}, keys ...interface{}) (interface{}, error) {
	defer UpdateStorageSQLMetric(OperationSelect, e.serverAddress, time.Now())

	entity, err := e.executor.Get(i, keys...)

	if err != nil {
		err = errors.New("error get data in DB, error: '" + err.Error() + "'")
	}

	return entity, err
}

func (e *SQLExecutor) Insert(list ...interface{}) error {
	defer UpdateStorageSQLMetric(OperationInsert, e.serverAddress, time.Now())

	if err := e.executor.Insert(list...); err != nil {
		return errors.New("error inserting data into DB, error: '" + err.Error() + "'")
	}

	return nil
}

func (e *SQLExecutor) Update(list ...interface{}) (int64, error) {
	defer UpdateStorageSQLMetric(OperationUpdate, e.serverAddress, time.Now())

	count, err := e.executor.Update(list...)

	if err != nil {
		err = errors.New("error updating data in DB, error: '" + err.Error() + "'")
	}

	return count, err
}

func (e *SQLExecutor) Delete(list ...interface{}) (int64, error) {
	defer UpdateStorageSQLMetric(OperationDelete, e.serverAddress, time.Now())

	count, err := e.executor.Delete(list...)

	if err != nil {
		err = errors.New("error deleting data in DB, error: '" + err.Error() + "'")
	}

	return count, err
}

func (e *SQLExecutor) Prepare(query string) (*sql.Stmt, error) {
	if transaction, ok := e.executor.(*gorp.Transaction); ok {
		return transaction.Prepare(query)
	} else if dbMap, ok := e.executor.(*gorp.DbMap); ok {
		return dbMap.Prepare(query)
	}

	return nil, errors.New("executor is not grop.Transaction or gorp.DbMap")
}

func (e *SQLExecutor) ExecByQuery(query string, args ...interface{}) (sql.Result, error) {
	defer UpdateStorageSQLMetric(OperationExec, e.serverAddress, time.Now())

	result, err := e.executor.Exec(query, args...)
	if err != nil {
		return result, errors.New("error executing DB query, query: '" + query + "', error: '" + err.Error() + "'")
	}
	return result, nil
}

func (e *SQLExecutor) Exec(query interface{}, args ...interface{}) (sql.Result, error) {
	if b, ok := query.(*sq.SelectBuilder); ok {
		return e.ExecSelect(b)
	}

	if b, ok := query.(*sq.InsertBuilder); ok {
		return e.ExecInsert(b)
	}

	if b, ok := query.(*sq.UpdateBuilder); ok {
		return e.ExecUpdate(b)
	}

	if b, ok := query.(*sq.DeleteBuilder); ok {
		return e.ExecDelete(b)
	}

	if b, ok := query.(*sq.CaseBuilder); ok {
		return e.ExecCase(b)
	}

	return nil, errors.New("could not prepare SQL query")
}

func (e *SQLExecutor) ExecSelect(builder *sq.SelectBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.New("could not prepare SQL query, error: '" + err.Error() + "'")
	}
	return e.ExecByQuery(query, args...)
}

func (e *SQLExecutor) ExecInsert(builder *sq.InsertBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.New("could not prepare SQL query, error: '" + err.Error() + "'")
	}
	return e.ExecByQuery(query, args...)
}

func (e *SQLExecutor) ExecUpdate(builder *sq.UpdateBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.New("could not prepare SQL query, error: '" + err.Error() + "'")
	}
	return e.ExecByQuery(query, args...)
}

func (e *SQLExecutor) ExecDelete(builder *sq.DeleteBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.New("could not prepare SQL query, error: '" + err.Error() + "'")
	}
	return e.ExecByQuery(query, args...)
}

func (e *SQLExecutor) ExecCase(builder *sq.CaseBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.New("could not prepare SQL query, error: '" + err.Error() + "'")
	}
	return e.ExecByQuery(query, args...)
}
