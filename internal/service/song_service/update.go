package song_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) Update(tx *sqlx.Tx, songId int, song models.Song) (updatedSong models.Song, err error) {
	exists, err := s.SongRepo.IsExists(tx, songId)
	if err != nil {
		return models.Song{}, err
	}
	if !exists {
		return models.Song{}, errors.NotFound{Resource: fmt.Sprintf("song with songId=%d in database", songId)}
	}

	err = s.SongRepo.Update(tx, songId, song)
	if err != nil {
		return models.Song{}, err
	}

	updatedSong, err = s.SongRepo.Read(tx, songId)
	if err != nil {
		return models.Song{}, err
	}

	return updatedSong, nil
}
