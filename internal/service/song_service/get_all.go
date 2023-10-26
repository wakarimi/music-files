package song_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

func (s *Service) GetAll(tx *sqlx.Tx) (songs []models.Song, err error) {
	songs, err = s.SongRepo.ReadAll(tx)
	if err != nil {
		return make([]models.Song, 0), err
	}

	return songs, nil
}
