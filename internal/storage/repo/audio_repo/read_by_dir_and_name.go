package audio_repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/audio"
)

func (r Repository) ReadByDirAndName(tx *sqlx.Tx, dirID int, name string) (audioFile audio.Audio, err error) {
	log.Debug().Int("dirId", dirID).Str("name", name).Msg("Reading audio file from database")

	query := `
		SELECT *
		FROM audios
		WHERE dir_id = :dir_id
			AND filename = :name
	`
	args := map[string]interface{}{
		"dir_id": dirID,
		"name":   name,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Str("name", name).Str("query", query).Msg("Failed to execute query to read audio file")
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
			log.Error().Err(err).Int("dirId", dirID).Str("name", name).Msg("Failed to get read result")
			return audio.Audio{}, err
		}
	} else {
		err := fmt.Errorf("no audio file found with dir_id: %d and name: %s", dirID, name)
		log.Error().Err(err).Int("dirId", dirID).Str("name", name).Msg("Audio file not found")
		return audio.Audio{}, err
	}

	log.Debug().Interface("audioFile", audioFile).Msg("Audio file read successfully")
	return audioFile, nil
}
