package audio_file_repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model"
)

func (r Repository) ReadAllBySha256(tx *sqlx.Tx, sha256 string) (audioFiles []model.AudioFile, err error) {
	log.Debug().Str("sha256", sha256).Msg("Reading audio files by sha256 from database")

	query := `
		SELECT * 
		FROM audio_files
		WHERE sha_256 = :sha_256
	`
	args := map[string]interface{}{
		"sha_256": sha256,
	}
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		log.Error().Err(err).Str("sha256", sha256).Str("query", query).Msg("Failed to execute query to read audio files by sha256")
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
			log.Error().Err(err).Str("sha256", sha256).Msg("Failed to get read result")
			return nil, err
		}
		audioFiles = append(audioFiles, audioFile)
	}

	log.Debug().Str("sha256", sha256).Int("countOfAudioFilesWithSha256", len(audioFiles)).Msg("Audio files by sha256 read successfully")
	return audioFiles, nil
}
