package song_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

func (s *Service) SearchBySha256(tx *sqlx.Tx, sha256 string) (songs []models.Song, err error) {
	songs, err = s.SongRepo.ReadAllBySha256(tx, sha256)
	if err != nil {
		return make([]models.Song, 0), err
	}

	return songs, nil
}
