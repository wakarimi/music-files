package audio_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) Delete(tx *sqlx.Tx, audioFileID int) (err error) {
	log.Debug().Int("audioFileId", audioFileID).Msg("Deleting audio file from database")

	query := `
		DELETE FROM audios
		WHERE id = :id
	`
	args := map[string]interface{}{
		"id": audioFileID,
	}
	_, err = tx.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("audioFileId", audioFileID).Msg("Failed to delete audio file from database")
		return err
	}

	log.Debug().Int("audioFileId", audioFileID).Msg("Audio file deleted from database successfully")
	return nil
}
