package audio_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) DeleteByDir(tx *sqlx.Tx, dirID int) (err error) {
	log.Debug().Int("dirId", dirID).Msg("Deleting audios by directory")

	query := `
		DELETE FROM audios
		WHERE dir_id = :dirID
	`
	args := map[string]interface{}{
		"dirID": dirID,
	}

	_, err = tx.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Msg("Failed to execute query to delete audios by directory")
		return err
	}

	log.Debug().Int("dirId", dirID).Msg("Audios deleted successfully by directory")
	return nil
}
