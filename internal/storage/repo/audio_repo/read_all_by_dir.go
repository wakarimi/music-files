package audio_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model/audio"
)

func (r Repository) ReadAllByDir(tx *sqlx.Tx, dirID int) (audioFiles []audio.Audio, err error) {
	log.Debug().Int("dirId", dirID).Msg("Reading audio files by directory from database")

	query := `
		SELECT * 
		FROM audios
		WHERE dir_id = :dir_id
	`
	args := map[string]interface{}{
		"dir_id": dirID,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirID).Str("query", query).Msg("Failed to execute query to read audio files by dirId")
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close rows")
		}
	}(rows)

	for rows.Next() {
		var audioFile audio.Audio
		if err = rows.StructScan(&audioFile); err != nil {
			log.Error().Err(err).Int("dirId", dirID).Msg("Failed to get read result")
			return nil, err
		}
		audioFiles = append(audioFiles, audioFile)
	}

	log.Debug().Int("dirId", dirID).Int("countOfAudioFilesInDir", len(audioFiles)).Msg("Audio files by dirId read successfully")
	return audioFiles, nil
}
