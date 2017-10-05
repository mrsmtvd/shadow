package database

import (
	"database/sql"

	"github.com/go-gorp/gorp"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

type Storage interface {
	GetDialect() string
	CreateTablesIfNotExists() error
	SetTypeConverter(converter gorp.TypeConverter)
	AddTableWithName(i interface{}, name string) *gorp.TableMap
	Begin() (Storage, error)
	Commit() error
	Rollback() error
	SelectByQuery(i interface{}, query string, args ...interface{}) ([]interface{}, error)
	Select(i interface{}, builder *sq.SelectBuilder) ([]interface{}, error)
	SelectOneByQuery(holder interface{}, query string, args ...interface{}) error
	SelectOne(holder interface{}, builder *sq.SelectBuilder) error
	SelectIntByQuery(query string, args ...interface{}) (int64, error)
	SelectInt(builder *sq.SelectBuilder) (int64, error)
	SelectNullIntByQuery(query string, args ...interface{}) (sql.NullInt64, error)
	SelectNullInt(builder *sq.SelectBuilder) (sql.NullInt64, error)
	SelectFloatByQuery(query string, args ...interface{}) (float64, error)
	SelectFloat(builder *sq.SelectBuilder) (float64, error)
	SelectNullFloatByQuery(query string, args ...interface{}) (sql.NullFloat64, error)
	SelectNullFloat(builder *sq.SelectBuilder) (sql.NullFloat64, error)
	SelectStrByQuery(query string, args ...interface{}) (string, error)
	SelectStr(builder *sq.SelectBuilder) (string, error)
	SelectNullStrByQuery(query string, args ...interface{}) (sql.NullString, error)
	SelectNullStr(builder *sq.SelectBuilder) (sql.NullString, error)
	Get(i interface{}, keys ...interface{}) (interface{}, error)
	Insert(list ...interface{}) error
	Update(list ...interface{}) (int64, error)
	Delete(list ...interface{}) (int64, error)
	Prepare(query string) (*sql.Stmt, error)
	ExecByQuery(query string, args ...interface{}) (sql.Result, error)
	Exec(query interface{}, args ...interface{}) (sql.Result, error)
	ExecSelect(builder *sq.SelectBuilder) (sql.Result, error)
	ExecInsert(builder *sq.InsertBuilder) (sql.Result, error)
	ExecUpdate(builder *sq.UpdateBuilder) (sql.Result, error)
	ExecDelete(builder *sq.DeleteBuilder) (sql.Result, error)
	ExecCase(builder *sq.CaseBuilder) (sql.Result, error)
}
