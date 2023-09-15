package track_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (s *Service) ReadAll(tx *sqlx.Tx) (tracks []models.Track, err error) {
	tracks, err = s.TrackRepo.ReadAllTx(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch all tracks")
		return make([]models.Track, 0), err
	}

	return tracks, nil
}
