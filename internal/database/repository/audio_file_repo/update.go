package audio_file_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model"
)

func (r Repository) Update(tx *sqlx.Tx, audioFileId int, audioFile model.AudioFile) (err error) {
	log.Debug().Int("audioFileId", audioFileId).Interface("audioFile", audioFile).Msg("Updating audio file")

	query := `
		UPDATE audio_files
		SET dir_id = :dir_id, filename = :filename, extension = :extension, size_byte = :size_byte,
		    duration_ms = :duration_ms, bitrate_kbps = :bitrate_kbps, sample_rate_hz = :sample_rate_hz,
		    channels_n = :channels_n, sha_256 = :sha_256, last_content_update = CURRENT_TIMESTAMP
		WHERE audio_file_id = :audio_file_id
	`

	audioFile.AudioFileId = audioFileId
	_, err = tx.NamedExec(query, audioFile)

	if err != nil {
		log.Error().Err(err).Int("audioFileId", audioFileId).Str("query", query).Msg("Failed to execute query to update audio file")
		return err
	}

	log.Debug().Int("audioFileId", audioFileId).Msg("Audio file updated successfully")
	return nil
}
