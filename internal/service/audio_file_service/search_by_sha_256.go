package audio_file_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

func (s *Service) SearchBySha256(tx *sqlx.Tx, sha256 string) (audioFiles []models.AudioFile, err error) {
	audioFiles, err = s.AudioFileRepo.ReadAllBySha256(tx, sha256)
	if err != nil {
		return make([]models.AudioFile, 0), err
	}

	return audioFiles, nil
}
