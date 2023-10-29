package audio_file_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) Update(tx *sqlx.Tx, audioFileId int, audioFile models.AudioFile) (updatedAudioFile models.AudioFile, err error) {
	exists, err := s.AudioFileRepo.IsExists(tx, audioFileId)
	if err != nil {
		return models.AudioFile{}, err
	}
	if !exists {
		return models.AudioFile{}, errors.NotFound{Resource: fmt.Sprintf("audioFile with audioFileId=%d in database", audioFileId)}
	}

	err = s.AudioFileRepo.Update(tx, audioFileId, audioFile)
	if err != nil {
		return models.AudioFile{}, err
	}

	updatedAudioFile, err = s.AudioFileRepo.Read(tx, audioFileId)
	if err != nil {
		return models.AudioFile{}, err
	}

	return updatedAudioFile, nil
}
