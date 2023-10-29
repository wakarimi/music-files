package audio_file_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) GetAudioFile(tx *sqlx.Tx, audioFileId int) (audioFile models.AudioFile, err error) {
	log.Debug().Int("audioFileId", audioFileId).Msg("Getting audioFile")

	exists, err := s.AudioFileRepo.IsExists(tx, audioFileId)
	if err != nil {
		return models.AudioFile{}, err
	}
	if !exists {
		err = errors.NotFound{Resource: fmt.Sprintf("audioFile with id=%d", audioFileId)}
		return models.AudioFile{}, err
	}

	audioFile, err = s.AudioFileRepo.Read(tx, audioFileId)
	if err != nil {
		return models.AudioFile{}, err
	}

	log.Debug().Int("audioFileId", audioFileId).Msg("AudioFile got successfully")
	return audioFile, nil
}
