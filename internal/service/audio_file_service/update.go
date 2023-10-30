package audio_file_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) Update(tx *sqlx.Tx, audioFileId int, audioFile models.AudioFile) (updatedAudioFile models.AudioFile, err error) {
	log.Debug().Int("audioFileId", audioFileId).Msg("Updating audio file")

	exists, err := s.AudioFileRepo.IsExists(tx, audioFileId)
	if err != nil {
		log.Debug().Int("audioFileId", audioFileId).Msg("Failed to check audio file existence")
		return models.AudioFile{}, err
	}
	if !exists {
		log.Error().Int("audioFileId", audioFileId).Msg("Audio file not found")
		return models.AudioFile{}, errors.NotFound{Resource: fmt.Sprintf("audioFile with audioFileId=%d in database", audioFileId)}
	}

	err = s.AudioFileRepo.Update(tx, audioFileId, audioFile)
	if err != nil {
		log.Debug().Int("audioFileId", audioFileId).Msg("Failed to update audio file")
		return models.AudioFile{}, err
	}

	updatedAudioFile, err = s.AudioFileRepo.Read(tx, audioFileId)
	if err != nil {
		log.Debug().Int("audioFileId", audioFileId).Msg("Failed to read updated audio file")
		return models.AudioFile{}, err
	}

	log.Debug().Int("audioFileId", audioFileId).Msg("Audio file updated successfully")
	return updatedAudioFile, nil
}
