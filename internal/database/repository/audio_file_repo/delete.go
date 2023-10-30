package audio_file_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r Repository) Delete(tx *sqlx.Tx, audioFileId int) (err error) {
	log.Debug().Int("audioFileId", audioFileId).Msg("Deleting audio file from database")

	query := `
		DELETE FROM audio_files
		WHERE audio_file_id = :audio_file_id
	`
	args := map[string]interface{}{
		"audio_file_id": audioFileId,
	}
	_, err = tx.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Int("audioFileId", audioFileId).Msg("Failed to delete audio file from database")
		return err
	}

	log.Debug().Int("audioFileId", audioFileId).Msg("Audio file deleted from database successfully")
	return nil
}
