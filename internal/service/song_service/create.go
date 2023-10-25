package song_service

import (
	"github.com/jmoiron/sqlx"
	"music-files/internal/models"
)

func (s *Service) Create(tx *sqlx.Tx, song models.Song) (createdSong models.Song, err error) {
	songId, err := s.SongRepo.Create(tx, song)
	if err != nil {
		return models.Song{}, err
	}

	createdSong, err = s.SongRepo.Read(tx, songId)
	if err != nil {
		return models.Song{}, err
	}

	return createdSong, nil
}
