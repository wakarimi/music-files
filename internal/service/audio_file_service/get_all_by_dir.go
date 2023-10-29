package audio_file_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

func (s *Service) GetAllByDir(tx *sqlx.Tx, dirId int) (audioFiles []models.AudioFile, err error) {
	audioFiles, err = s.AudioFileRepo.ReadAllByDir(tx, dirId)
	if err != nil {
		return make([]models.AudioFile, 0), err
	}

	return audioFiles, nil
}
