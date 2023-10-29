package audio_file_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

func (s *Service) GetAll(tx *sqlx.Tx) (audioFiles []models.AudioFile, err error) {
	audioFiles, err = s.AudioFileRepo.ReadAll(tx)
	if err != nil {
		return make([]models.AudioFile, 0), err
	}

	return audioFiles, nil
}
