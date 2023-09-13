package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Queryer interface {
	Get(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	Queryx(query string, arg ...interface{}) (*sqlx.Rows, error)
	Select(dest interface{}, query string, args ...interface{}) error
}
