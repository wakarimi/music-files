package audio_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) IsExistsByDirAndName(tx *sqlx.Tx, dirID int, name string) (exists bool, err error) {
	log.Debug().Int("dirId", dirID).Str("name", name).Msg("Checking for the existence of a audio file in the database by dirId and name")

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM audios
			WHERE dir_id = :dir_id
				AND filename = :name
		)
	`
	args := map[string]interface{}{
		"dir_id": dirID,
		"name":   name,
	}
	row, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Str("name", name).Str("query", query).Msg("Failed to execute query to check existence in database")
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
			log.Error().Err(err).Int("dirId", dirID).Str("name", name).Msg("Failed to get existence check results")
			return false, err
		}
	}

	log.Debug().Int("dirId", dirID).Str("name", name).Bool("exists", exists).Msg("The existence of the audio file was checked successfully")
	return exists, nil
}
