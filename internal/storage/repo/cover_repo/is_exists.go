package cover_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsExists(tx *sqlx.Tx, coverID int) (exists bool, err error) {
	log.Debug().Int("coverId", coverID).Msg("Checking for the existence of a cover file in the database")

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM covers
			WHERE id = :id
		)
	`
	args := map[string]interface{}{
		"id": coverID,
	}
	row, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("coverId", coverID).Str("query", query).Msg("Failed to execute query to check existence in database")
		return false, err
	}
	defer func(row *sqlx.Rows) {
		err := row.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close row")
		}
	}(row)
	if row.Next() {
		if err = row.Scan(&exists); err != nil {
			log.Error().Err(err).Int("coverId", coverID).Msg("Failed to get existence check results")
			return false, err
		}
	}

	log.Debug().Int("coverId", coverID).Bool("exists", exists).Msg("The existence of the cover file was checked successfully")
	return exists, nil
}
