package audio_file_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) GetAudioFile(tx *sqlx.Tx, audioFileId int) (audioFile models.AudioFile, err error) {
	log.Debug().Int("audioFileId", audioFileId).Msg("Getting audio file")

	exists, err := s.AudioFileRepo.IsExists(tx, audioFileId)
	if err != nil {
		log.Error().Int("audioFileId", audioFileId).Msg("Failed to check audio file existence")
		return models.AudioFile{}, err
	}
	if !exists {
		log.Error().Int("audioFileId", audioFileId).Msg("Audio file not found")
		return models.AudioFile{}, errors.NotFound{Resource: fmt.Sprintf("audio_file with id=%d", audioFileId)}
	}

	audioFile, err = s.AudioFileRepo.Read(tx, audioFileId)
	if err != nil {
		log.Error().Err(err).Int("audioFileId", audioFileId).Msg("Failed to fetch audio file")
		return models.AudioFile{}, err
	}

	log.Debug().Int("audioFileId", audioFileId).Msg("Audio file got successfully")
	return audioFile, nil
}
