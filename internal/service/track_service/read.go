package track_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/models"
)

func (s *Service) Read(tx *sqlx.Tx, trackId int) (track models.Track, err error) {
	track, err = s.TrackRepo.ReadTx(tx, trackId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read track")
		return models.Track{}, err
	}

	return track, nil
}
