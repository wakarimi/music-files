package song_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"music-files/internal/models"
)

func (s *Service) GetSong(tx *sqlx.Tx, songId int) (song models.Song, err error) {
	log.Debug().Int("songId", songId).Msg("Getting directory")

	exists, err := s.SongRepo.IsExists(tx, songId)
	if err != nil {
		return models.Song{}, err
	}
	if !exists {
		err = errors.NotFound{Resource: fmt.Sprintf("song with id=%d", songId)}
		return models.Song{}, err
	}

	song, err = s.SongRepo.Read(tx, songId)
	if err != nil {
		return models.Song{}, err
	}

	log.Debug().Int("songId", songId).Msg("Song got successfully")
	return song, nil
}
