package dir_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r *Repository) IsExists(tx *sqlx.Tx, dirId int) (exists bool, err error) {
	log.Debug().Int("dirId", dirId).
		Msg("Checking if directory exists in database")

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM directories
			WHERE dir_id = :dir_id
		)
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	row, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).
			Msg("Failed to execute query to check directory existence")
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
			log.Error().Err(err).Int("dirId", dirId).
				Msg("Failed to scan result of directory existence check")
			return false, err
		}
	}

	log.Debug().Int("dirId", dirId).Bool("exists", exists)
	return exists, nil
}
