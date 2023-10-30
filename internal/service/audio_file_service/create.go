package audio_file_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (s *Service) Create(tx *sqlx.Tx, audioFile models.AudioFile) (createdAudioFile models.AudioFile, err error) {
	log.Debug().Interface("audioFile", audioFile).Msg("Creating new audio file")

	audioFileId, err := s.AudioFileRepo.Create(tx, audioFile)
	if err != nil {
		log.Error().Err(err).Interface("audioFile", audioFile).Msg("Failed to create new audio file")
		return models.AudioFile{}, err
	}

	createdAudioFile, err = s.AudioFileRepo.Read(tx, audioFileId)
	if err != nil {
		log.Error().Err(err).Interface("audioFile", audioFile).Msg("Failed to read created audio file")
		return models.AudioFile{}, err
	}

	log.Debug().Interface("audioFile", audioFile).Msg("Audio file created successfully")
	return createdAudioFile, nil
}
