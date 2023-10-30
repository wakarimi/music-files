package audio_file_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (s *Service) GetAll(tx *sqlx.Tx) (audioFiles []models.AudioFile, err error) {
	log.Debug().Msg("Fetching all audio files")

	audioFiles, err = s.AudioFileRepo.ReadAll(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch all audio files")
		return make([]models.AudioFile, 0), err
	}

	log.Debug().Int("countOfAudioFiles", len(audioFiles)).Msg("All audio files fetched successfully")
	return audioFiles, nil
}
