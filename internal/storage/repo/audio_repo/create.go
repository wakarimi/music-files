package audio_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/audio"
)

func (r Repository) Create(tx *sqlx.Tx, audioFile audio.Audio) (audioFileID int, err error) {
	log.Debug().Interface("audioFile", audioFile).Msg("Creating new audio file in database")

	query := `
		INSERT INTO audios(dir_id, filename, extension, size_byte, duration_ms, bitrate_kbps, sample_rate_hz, channels_n, sha_256, last_content_update)
		VALUES (:dir_id, :filename, :extension, :size_byte, :duration_ms, :bitrate_kbps, :sample_rate_hz, :channels_n, :sha_256, CURRENT_TIMESTAMP)
		RETURNING id
	`
	rows, err := tx.NamedQuery(query, audioFile)
	if err != nil {
		log.Error().Err(err).Interface("audioFile", audioFile).Str("query", query).Msg("Failed to create audio file in database")
		return 0, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}(rows)

	if rows.Next() {
		if err := rows.Scan(&audioFileID); err != nil {
			log.Error().Err(err).Msg("Failed to scan coverId of created audio file")
			return 0, err
		}
	} else {
		err := fmt.Errorf("no id returned after audio_file insert")
		log.Error().Err(err).Interface("audioFile", audioFile).Msg("No id returned after audio file insert")
		return 0, err
	}

	log.Debug().Int("audioFileId", audioFileID).Interface("audioFile", audioFile).Msg("New audio file in database created successfully")
	return audioFileID, nil
}
