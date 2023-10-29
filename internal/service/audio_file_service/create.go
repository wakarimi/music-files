package audio_file_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

func (s *Service) Create(tx *sqlx.Tx, audioFile models.AudioFile) (createdAudioFile models.AudioFile, err error) {
	audioFileId, err := s.AudioFileRepo.Create(tx, audioFile)
	if err != nil {
		return models.AudioFile{}, err
	}

	createdAudioFile, err = s.AudioFileRepo.Read(tx, audioFileId)
	if err != nil {
		return models.AudioFile{}, err
	}

	return createdAudioFile, nil
}
