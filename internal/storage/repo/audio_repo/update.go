package audio_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/audio"
)

func (r Repository) Update(tx *sqlx.Tx, audioFileID int, audioFile audio.Audio) (err error) {
	log.Debug().Int("audioFileId", audioFileID).Interface("audioFile", audioFile).Msg("Updating audio file")

	query := `
		UPDATE audios
		SET dir_id = :dir_id, filename = :filename, extension = :extension, size_byte = :size_byte,
		    duration_ms = :duration_ms, bitrate_kbps = :bitrate_kbps, sample_rate_hz = :sample_rate_hz,
		    channels_n = :channels_n, sha_256 = :sha_256, last_content_update = CURRENT_TIMESTAMP
		WHERE id = :id
	`

	audioFile.ID = audioFileID
	_, err = tx.NamedExec(query, audioFile)

	if err != nil {
		log.Error().Err(err).Int("audioFileId", audioFileID).Str("query", query).Msg("Failed to execute query to update audio file")
		return err
	}

	log.Debug().Int("audioFileId", audioFileID).Msg("Audio file updated successfully")
	return nil
}
