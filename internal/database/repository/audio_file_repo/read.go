package audio_file_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model"
)

func (r Repository) Read(tx *sqlx.Tx, audioFileId int) (audioFile model.AudioFile, err error) {
	log.Debug().Int("audioFileId", audioFileId).Msg("Reading audio file from database")

	query := `
		SELECT *
		FROM audio_files
		WHERE audio_file_id = :audio_file_id
	`
	args := map[string]interface{}{
		"audio_file_id": audioFileId,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("audioFileId", audioFileId).Str("query", query).Msg("Failed to execute query to read audio file")
		return model.AudioFile{}, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}(rows)
	if rows.Next() {
		if err = rows.StructScan(&audioFile); err != nil {
			log.Error().Err(err).Int("audioFileId", audioFileId).Msg("Failed to get read result")
			return model.AudioFile{}, err
		}
	} else {
		err := fmt.Errorf("No audio file found with audio_file_id: %d", audioFileId)
		log.Error().Err(err).Int("audioFileId", audioFileId).Msg("Audio file not found")
		return model.AudioFile{}, err
	}

	log.Debug().Interface("audioFileId", audioFileId).Msg("Audio file read successfully")
	return audioFile, nil
}
