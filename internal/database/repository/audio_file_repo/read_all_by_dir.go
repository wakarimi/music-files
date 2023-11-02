package audio_file_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model"
)

func (r Repository) ReadAllByDir(tx *sqlx.Tx, dirId int) (audioFiles []model.AudioFile, err error) {
	log.Debug().Int("dirId", dirId).Msg("Reading audio files by directory from database")

	query := `
		SELECT * 
		FROM audio_files
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirId,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Str("query", query).Msg("Failed to execute query to read audio files by dirId")
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}(rows)

	for rows.Next() {
		var audioFile model.AudioFile
		if err = rows.StructScan(&audioFile); err != nil {
			log.Error().Err(err).Int("dirId", dirId).Msg("Failed to get read result")
			return nil, err
		}
		audioFiles = append(audioFiles, audioFile)
	}

	log.Debug().Int("dirId", dirId).Int("countOfAudioFilesInDir", len(audioFiles)).Msg("Audio files by dirId read successfully")
	return audioFiles, nil
}
