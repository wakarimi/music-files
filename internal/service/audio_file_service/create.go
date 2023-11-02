package audio_file_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/model"
)

func (s *Service) Create(tx *sqlx.Tx, audioFile model.AudioFile) (createdAudioFile model.AudioFile, err error) {
	log.Debug().Interface("audioFile", audioFile).Msg("Creating new audio file")

	audioFileId, err := s.AudioFileRepo.Create(tx, audioFile)
	if err != nil {
		log.Error().Err(err).Interface("audioFile", audioFile).Msg("Failed to create new audio file")
		return model.AudioFile{}, err
	}

	createdAudioFile, err = s.AudioFileRepo.Read(tx, audioFileId)
	if err != nil {
		log.Error().Err(err).Interface("audioFile", audioFile).Msg("Failed to read created audio file")
		return model.AudioFile{}, err
	}

	log.Debug().Interface("audioFile", audioFile).Msg("Audio file created successfully")
	return createdAudioFile, nil
}
