package service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type TransactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db sqlx.DB) (txManager TransactionManager) {
	return TransactionManager{
		db: &db,
	}
}

func (tm *TransactionManager) begin() (tx *sqlx.Tx, err error) {
	log.Debug().Msg("Starting a new transaction")
	tx, err = tm.db.Beginx()
	return tx, err
}

func (tm *TransactionManager) WithTransaction(do func(tx *sqlx.Tx) (err error)) (err error) {
	tx, err := tm.begin()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start a transaction")
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Error().Msg("Recovered from panic, rolling back transaction")
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Error().Err(err).Msg("An error occurred, rolling back transaction")
			tx.Rollback()
		} else {
			log.Debug().Msg("Committing transaction")
			err = tx.Commit()
		}
	}()

	err = do(tx)
	return err
}
