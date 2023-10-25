package song_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

func (s *Service) GetAllByDir(tx *sqlx.Tx, dirId int) (songs []models.Song, err error) {
	songs, err = s.SongRepo.ReadAllByDir(tx, dirId)
	if err != nil {
		return make([]models.Song, 0), err
	}

	return songs, nil
}
