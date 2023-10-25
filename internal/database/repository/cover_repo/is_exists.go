package cover_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsExists(tx *sqlx.Tx, coverId int) (exists bool, err error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM covers
			WHERE cover_id = :cover_id
		)
	`
	args := map[string]interface{}{
		"cover_id": coverId,
	}
	row, err := tx.NamedQuery(query, args)
	if err != nil {
		return false, err
	}
	defer func(row *sqlx.Rows) {
		err := row.Close()
		if err != nil {
			log.Error().Err(err)
		}
	}(row)
	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			return false, err
		}
	}

	return exists, nil
}
