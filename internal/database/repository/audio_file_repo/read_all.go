package audio_file_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model"
)

func (r *Repository) ReadAll(tx *sqlx.Tx) (audioFiles []model.AudioFile, err error) {
	log.Debug().Msg("Reading all audio files")

	query := `
		SELECT * 
		FROM audio_files
	`
	err = tx.Select(&audioFiles, query)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("Failed to execute query to read audio files")
		return nil, err
	}

	log.Debug().Int("countOfAudioFiles", len(audioFiles)).Msg("All audio files fetched successfully")
	return audioFiles, nil
}
