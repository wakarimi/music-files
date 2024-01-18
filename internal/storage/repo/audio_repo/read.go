package audio_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/audio"
)

func (r Repository) Read(tx *sqlx.Tx, audioID int) (audioFile audio.Audio, err error) {
	log.Debug().Int("audioId", audioID).Msg("Reading audio file from database")

	query := `
		SELECT *
		FROM audios
		WHERE id = :id
	`
	args := map[string]interface{}{
		"id": audioID,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("audioId", audioID).Str("query", query).Msg("Failed to execute query to read audio file")
		return audio.Audio{}, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}(rows)
	if rows.Next() {
		if err = rows.StructScan(&audioFile); err != nil {
			log.Error().Err(err).Int("audioId", audioID).Msg("Failed to get read result")
			return audio.Audio{}, err
		}
	} else {
		err := fmt.Errorf("No audio file found with id: %d", audioID)
		log.Error().Err(err).Int("audioId", audioID).Msg("Audio file not found")
		return audio.Audio{}, err
	}

	log.Debug().Interface("audioId", audioID).Msg("Audio file read successfully")
	return audioFile, nil
}
