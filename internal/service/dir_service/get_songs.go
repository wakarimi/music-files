package dir_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (s *Service) GetSongs(tx *sqlx.Tx, dirId int) (songs []models.Song, err error) {
	log.Debug().Msg("Getting songs in directory")

	songs, err = s.SongService.GetAllByDir(tx, dirId)
	if err != nil {
		log.Error().Err(err).Int("dirId", dirId).Msg("Failed to get songs in directory")
		return make([]models.Song, 0), err
	}

	log.Debug().Int("countOfSongs", len(songs)).Msg("Song in directory got successfully")
	return songs, nil
}
