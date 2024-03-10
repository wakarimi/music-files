package txsql

import "github.com/jmoiron/sqlx"

type SQLTransaction struct {
	tx *sqlx.Tx
}

func New(tx *sqlx.Tx) *SQLTransaction {
	return &SQLTransaction{
		tx: tx,
	}
}

func (t *SQLTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *SQLTransaction) Rollback() error {
	return t.tx.Commit()
}
