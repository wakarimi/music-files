package audio_file_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsExists(tx *sqlx.Tx, audioFileId int) (exists bool, err error) {
	log.Debug().Int("audioFileId", audioFileId).Msg("Checking for the existence of a audio file in the database")

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM audio_files
			WHERE audio_file_id = :audio_file_id
		)
	`
	args := map[string]interface{}{
		"audio_file_id": audioFileId,
	}
	row, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("audioFileId", audioFileId).Str("query", query).Msg("Failed to execute query to check existence in database")
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
			log.Error().Err(err).Int("audioFileId", audioFileId).Msg("Failed to get existence check results")
			return false, err
		}
	}

	log.Debug().Int("audioFileId", audioFileId).Bool("exists", exists).Msg("The existence of the audio file was checked successfully")
	return exists, nil
}
