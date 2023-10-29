package audio_file_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (r Repository) ReadByDirAndName(tx *sqlx.Tx, dirId int, name string) (audioFile models.AudioFile, err error) {
	log.Debug().Int("dirId", dirId).Str("name", name).Msg("Reading audio file from database")

	query := `
		SELECT *
		FROM audio_files
		WHERE dir_id = :dir_id
			AND filename = :name
	`
	args := map[string]interface{}{
		"dir_id": dirId,
		"name":   name,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Str("name", name).Str("query", query).Msg("Failed to execute query to read audio file")
		return models.AudioFile{}, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}(rows)
	if rows.Next() {
		if err = rows.StructScan(&audioFile); err != nil {
			log.Error().Err(err).Int("dirId", dirId).Str("name", name).Msg("Failed to get read result")
			return models.AudioFile{}, err
		}
	} else {
		err := fmt.Errorf("no audio file found with dir_id: %d and name: %s", dirId, name)
		log.Error().Err(err).Int("dirId", dirId).Str("name", name).Msg("Audio file not found")
		return models.AudioFile{}, err
	}

	log.Debug().Interface("audioFile", audioFile).Msg("Audio file read successfully")
	return audioFile, nil
}
