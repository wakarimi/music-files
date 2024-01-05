package service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Transactor struct {
	db *sqlx.DB
}

func NewTransactor(db sqlx.DB) (transactor *Transactor) {
	return &Transactor{
		db: &db,
	}
}

func (t *Transactor) begin() (tx *sqlx.Tx, err error) {
	log.Debug().Msg("Starting a new transaction")
	tx, err = t.db.Beginx()
	return tx, err
}

func (t *Transactor) WithTransaction(do func(tx *sqlx.Tx) (err error)) (err error) {
	tx, err := t.begin()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start a transaction")
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Error().Msg("Recovered from panic, rolling back transaction")
			err := tx.Rollback()
			if err != nil {
				log.Error().Err(err).Msg("Failed to rollback transaction")
				return
			}
			panic(p)
		} else if err != nil {
			log.Error().Err(err).Msg("An error occurred, rolling back transaction")
			err := tx.Rollback()
			if err != nil {
				log.Error().Err(err).Msg("Failed to rollback transaction")
				return
			}
		} else {
			log.Debug().Msg("Committing transaction")
			err = tx.Commit()
		}
	}()

	err = do(tx)
	return err
}
